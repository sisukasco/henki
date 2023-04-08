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

type ResetPasswordInfo struct {
	FirstName string
	Link      string
}

func (usvc *UserService) generatePasswordResetToken(ctx context.Context, userID string) (string, error) {
	resetToken := stringid.RandString(22)
	for i := 0; i < 100; i++ {
		exists, err := usvc.svc.DB.Q.DoesPasswordResetTokenExist(ctx, resetToken)
		if err != nil {
			log.Printf("Error checking for password reset token: %v", err)
			return "", err
		}
		if !exists {
			break
		}
		resetToken = stringid.RandString(22)
	}

	err := usvc.svc.DB.Q.UpdatePasswordResetToken(ctx,
		db.UpdatePasswordResetTokenParams{ID: userID, RecoveryToken: resetToken})

	if err != nil {
		log.Printf("Error while sending email Error updating confirm code %v", err)
		return "", err
	}

	return resetToken, nil
}

//sendPasswordResetEmail is called from task that is called from job queue
//Composes and sends the password reset email. Does not send email if called
//more than once in an hour
func (usvc *UserService) sendPasswordResetEmail(ctx context.Context, email string) {

	user, err := usvc.svc.DB.Q.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Error while sending reset password link:  user not found %v", err)
		return
	}
	

	resetToken, err := usvc.generatePasswordResetToken(ctx, user.ID)

	site := usvc.svc.Konf.String("client.url")

	subj := usvc.svc.Konf.String("emails.reset-password.subject")
	resetPasswordEmail := usvc.svc.Konf.String("emails.reset-password.body")
	link := site + "/auth/reset/?token=" + resetToken
	e := ResetPasswordInfo{user.FirstName, link}
	eb, err := mustache.Render(resetPasswordEmail, &e)
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
	log.Println("reset password email sent")
}

func (usvc *UserService) InitResetPasswordRequest(ctx context.Context, email string) error {
	email = utils.CleanupString(email)

	user, err := usvc.svc.DB.Q.GetUserByEmail(ctx, email)
	if err != nil {
		return http_utils.BadRequestError("This email was not used to register with this service")
	}

	if user.RecoverySentAt.Valid {
		last := time.Now().Sub(user.RecoverySentAt.Time).Minutes()
		if last < 10 { //If email was sent in last 10 minutes, dont send again

			return http_utils.BadRequestError("Another password reset request in progress. Please wait 10 minutes before another attempt")
		}
	}

	usvc.PostPasswordResetEmail(email)

	return nil
}
func (usvc *UserService) PostPasswordResetEmail(email string) {
	usvc.wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		defer cancel()
		defer usvc.wg.Done()
		usvc.sendPasswordResetEmail(ctx, email)
	}()
}

func (usvc *UserService) ResetPassword(ctx context.Context, token string, passwd string) error {
	token = strings.TrimSpace(token)
	if len(token) <= 0 {
		return http_utils.BadRequestError("The token is empty!")
	}
	if len(token) > 50 {
		return http_utils.BadRequestError("Bad Token")
	}
	passwd = strings.TrimSpace(passwd)
	if len(passwd) <= 0 {
		return http_utils.BadRequestError("The Password is empty!")
	}
	if len(passwd) > 250 {
		return http_utils.BadRequestError("The Password is too big!")
	}

	rec, err := usvc.svc.DB.Q.GetUserFromRecoveryToken(ctx, token)
	if err != nil {
		return http_utils.BadRequestError("Invalid Token")
	}
	if !rec.RecoverySentAt.Valid {
		return http_utils.BadRequestError("The token is not valid.")
	}
	diff := time.Now().Sub(rec.RecoverySentAt.Time).Hours()
	if diff > 24 {
		return http_utils.BadRequestError("The recovery token has expired. Try resetting the password again.")
	}
	encrpswd, err := hashPassword(passwd)
	if err != nil {
		log.Printf("Error %v", err)
		return http_utils.InternalServerError("Error updating record").WithInternalError(err)
	}
	err = usvc.svc.DB.Q.UpdateRecoveryPassword(ctx, db.UpdateRecoveryPasswordParams{ID: rec.ID, EncryptedPassword: encrpswd})
	if err != nil {
		log.Printf("Error %v", err)
		return http_utils.InternalServerError("Error updating record").WithInternalError(err)
	}

	err = usvc.updateResetPasswordOnConfirmation(ctx, rec.ID, false)
	if err != nil {
		log.Printf("Error updating ResetPasswordOnConfirmation %v", err)
	}

	return nil
}

func (usvc *UserService) UpdatePassword(ctx context.Context, user *db.User, oldPassword string, newPassword string) error {

	oldPassword = strings.TrimSpace(oldPassword)
	newPassword = strings.TrimSpace(newPassword)

	if !user.Authenticate(oldPassword) {
		return http_utils.UnauthorizedError("Old password does not match")
	}
	if len(newPassword) > 250 {
		return http_utils.UnprocessableEntityError("Too long password ")
	}

	pw, err := hashPassword(newPassword)
	if err != nil {
		return http_utils.InternalServerError("Couldn't create password hash").WithInternalError(err)
	}

	err = usvc.svc.DB.Q.UpdatePassword(ctx, db.UpdatePasswordParams{user.ID, pw})

	if err != nil {
		return http_utils.InternalServerError("Couldn't update password").WithInternalError(err)
	}
	return nil
}
