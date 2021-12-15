package usersvc

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/commons/stringid"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"

	"github.com/cbroglie/mustache"
)

type ConfirmEmailInfo struct {
	FirstName string
	Link      string
}

func (usvc *UserService) SendEmailConfirmationRequest(ctx context.Context, userID string) error {

	user, err := usvc.svc.DB.Q.GetUser(ctx, userID)
	if err != nil {
		return http_utils.InternalServerError("Error fetching records").WithInternalError(err)
	}
	if user.ConfirmedAt.Valid {
		//already confirmed
		return nil
	}

	if user.ConfirmationSentAt.Valid {
		last := time.Now().Sub(user.ConfirmationSentAt.Time).Minutes()
		if last < 60 { //If email was sent in last 50 minutes, dont send again
			return nil
		}
	}

	usvc.PostConfirmationEmail(user.ID)

	return nil
}

func (usvc *UserService) PostConfirmationEmail(userID string) {
	usvc.wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	go func() {
		defer cancel()
		defer usvc.wg.Done()
		usvc.sendConfirmationEmail(ctx, userID)
	}()

}

type ConfirmEmailResponse struct {
	ResetPasswordToken string `json:"reset_token"`
	Status             string `json:"status"`
}

func (usvc *UserService) ConfirmUserEmail(ctx context.Context, code string) (*ConfirmEmailResponse, error) {

	code = strings.TrimSpace(code)

	uinfo, err := usvc.svc.DB.Q.GetUserByConfirmationToken(ctx, code)
	if err != nil || len(uinfo.ID) <= 2 {
		log.Printf(" ConfirmUserEmail err %v, uinfo %v", err, utils.ToJSONString(uinfo))
		return nil, http_utils.UnprocessableEntityError("Confirmation code does not match").WithInternalError(err)
	}

	err = usvc.svc.DB.Q.ConfirmUserEmail(ctx, code)
	if err != nil {
		return nil, http_utils.InternalServerError("Error while accessing").WithInternalError(err)
	}
	resp := &ConfirmEmailResponse{}
	if uinfo.Resetpassword {
		token, err := usvc.generatePasswordResetToken(ctx, uinfo.ID)
		if err != nil {
			return nil, err
		}

		resp.ResetPasswordToken = token

	}
	resp.Status = "ok"
	return resp, nil
}

//sendConfirmationEmail is called from the task (in the job queue)
// So this function is called asynchronously
func (usvc *UserService) sendConfirmationEmail(ctx context.Context, userID string) {

	log.Printf("sending confirmation email to %s ", userID)
	user, err := usvc.svc.DB.Q.GetUser(ctx, userID)
	if err != nil {
		log.Printf("Error while sending email user not found %v", err)
		return
	}

	confirmCode := stringid.RandString(12)
	for i := 0; i < 100; i++ {
		exists, err := usvc.svc.DB.Q.DoesConfirmationTokenExist(ctx, confirmCode)
		if !exists && err == nil {
			break
		}
		confirmCode = stringid.RandString(12)
	}

	err = usvc.svc.DB.Q.UpdateConfirmationToken(ctx,
		db.UpdateConfirmationTokenParams{ID: user.ID, ConfirmationToken: confirmCode})

	if err != nil {
		log.Printf("Error while sending email Error updating confirm code %v", err)
		return
	}

	site := usvc.svc.Konf.String("client.url")

	subj := usvc.svc.Konf.String("emails.signup-confirmation.subject")
	confirmationEmail := usvc.svc.Konf.String("emails.signup-confirmation.body")

	confLink := site + "/auth/confirm/?code=" + confirmCode
	e := ConfirmEmailInfo{user.FirstName, confLink}
	eb, err := mustache.Render(confirmationEmail, &e)
	if err != nil {
		log.Printf("Error while sending email- mustache rendering error %v", err)
		return
	}
	mailer := usvc.getMailer()
	err = mailer.SendEmail(user.Email, subj, eb)
	if err != nil {
		log.Printf("Error while sending email- SendEmail error %v", err)
		return
	}
	log.Println("confirmation email sent")
}
