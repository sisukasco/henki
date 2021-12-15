package usersvc

import (
	"context"
	"log"
	"strings"
	"time"

	emailer "github.com/sisukasco/commons/email"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/commons/stringid"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"

	"github.com/cbroglie/mustache"
	"github.com/pkg/errors"
)

type EmailUpdateInfo struct {
	FirstName string
	Link      string
}

func (usvc *UserService) sendEmailUpdateEmail(ctx context.Context, userID string,
	new_email string) {

	log.Printf("Beginning email update for user ID %s", userID)

	user, err := usvc.svc.DB.Q.GetUser(ctx, userID)
	if err != nil {
		log.Printf("Error while sending email change link:  user not found %v", err)
		return
	}
	token := stringid.RandString(22)
	for i := 0; i < 100; i++ {
		exists, err := usvc.svc.DB.Q.DoesEmailUpdateTokenExist(ctx, token)
		if !exists || err == nil {
			break
		}
		token = stringid.RandString(22)
	}

	err = usvc.svc.DB.Q.InitUpdateUserEmail(ctx,
		db.InitUpdateUserEmailParams{ID: user.ID, EmailChange: new_email, EmailChangeToken: token})

	if err != nil {
		log.Printf("Error while sending email change link: Error updating token code %v", err)
		return
	}
	site := usvc.svc.Konf.String("client.url")

	subj := usvc.svc.Konf.String("emails.email-update.subject")
	templateUpdateEmail := usvc.svc.Konf.String("emails.email-update.body")
	link := site + "/auth/email-update/?token=" + token
	e := EmailUpdateInfo{user.FirstName, link}
	eb, err := mustache.Render(templateUpdateEmail, &e)
	if err != nil {
		log.Printf("Error while sending email- mustache rendering error %v", err)
		return
	}
	mailer := usvc.getMailer()
	err = mailer.SendEmail(new_email, subj, eb)
	if err != nil {
		log.Printf("Error while sending email- SendEmail error %v", err)
		return
	}
	log.Println("email update confirmation email sent")
}

func (usvc *UserService) getMailerConfig() *emailer.EmailConfig {
	if usvc.svc.Konf.Exists("mailer.smtp") {
		ss := usvc.svc.Konf.StringMap("mailer.smtp")
		return emailer.NewSMTPConfig(usvc.svc.Konf.String("mailer.from"),
			ss["host"], ss["user"], ss["pass"])
	} else if usvc.svc.Konf.Exists("mailer.ses") {
		ss := usvc.svc.Konf.StringMap("mailer.ses")
		return emailer.NewSESConfig(usvc.svc.Konf.String("mailer.from"),
			ss["region"], ss["access_key"], ss["secret_key"])
	} else {
		return &emailer.EmailConfig{}
	}

}
func (usvc *UserService) getMailer() *emailer.Emailer {
	return emailer.NewEmailer(usvc.getMailerConfig())
}

func (usvc *UserService) UpdateProfileField(ctx context.Context, user *db.User,
	fieldName string, fieldValue string) error {

	var err error

	switch fieldName {
	case "first_name":
		err = usvc.svc.DB.Q.UpdateUserFirstName(ctx,
			db.UpdateUserFirstNameParams{user.ID, fieldValue})

	case "last_name":
		err = usvc.svc.DB.Q.UpdateUserLastName(ctx,
			db.UpdateUserLastNameParams{user.ID, fieldValue})

	case "email":
		email := utils.CleanupString(fieldValue)

		if len(email) <= 0 {
			return http_utils.UnprocessableEntityError("Email can't be empty")
		}

		exists, _ := usvc.svc.DB.Q.DoesUserExist(ctx, email)
		if exists {
			return http_utils.UnprocessableEntityError("A user with the same email address already exists.")
		}

		usvc.PostEmailUpdateConfirmation(user.ID, email)

	default:
		return errors.New("Unknown Field name " + fieldName)
	}

	return err
}

func (usvc *UserService) PostEmailUpdateConfirmation(userID string, email string) {
	usvc.wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		defer cancel()
		defer usvc.wg.Done()
		usvc.sendEmailUpdateEmail(ctx, userID, email)
	}()
}

func (usvc *UserService) CompleteEmailUpdate(ctx context.Context, token string) error {

	token = strings.TrimSpace(token)

	ur, err := usvc.svc.DB.Q.GetUserFromEmailUpdateToken(ctx, token)
	if err != nil {
		return http_utils.InternalServerError("Please check the token code").WithInternalError(err)
	}

	if len(ur.EmailChange) <= 0 {
		return http_utils.UnprocessableEntityError("Bad token %s", token)
	}

	diff := time.Now().Sub(ur.EmailChangeSentAt.Time).Hours()
	if diff > 48 {
		return http_utils.UnprocessableEntityError("This token got expired. Token gets expired after 48 hours")
	}

	err = usvc.svc.DB.Q.UpdateUserEmail(ctx, db.UpdateUserEmailParams{
		ur.ID, ur.EmailChange})

	if err != nil {
		return http_utils.InternalServerError("Couldn't get User Rec").WithInternalError(err)
	}
	return nil
}
