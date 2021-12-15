package usersvc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"
	"github.com/sisukasco/henki/pkg/external"
	"github.com/sisukasco/henki/pkg/service"
	dtesting "github.com/sisukasco/henki/pkg/testing"
	"github.com/sisukasco/henki/pkg/usersvc"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

var (
	ctx context.Context
	svc *service.Service
)

func TestMain(m *testing.M) {
	dtesting.InitService(m, func(s *service.Service) {
		svc = s
	}, func() {
	})

}

func TestSignup(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    faker.Internet().Email(),
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)

	if err != nil {
		t.Errorf("Error signingup new user %v ", err)
		return
	}

	userNew, err := usersvc.GetUserByEmail(ctx, svc.DB.Q, params.Email)
	if err != nil {
		t.Errorf("Error signingup Failed Getting the user back %v ", err)
		return
	}

	t.Logf("Signing up user - got user back %v", userNew.ID)
}

func TestSignupNoEmail(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    "  ",
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)
	if err == nil {
		t.Errorf("Signup allows empty email !")
	} else {
		t.Logf("Received error %v", err)
	}

}

func TestSignupTooBigEmail(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    faker.RandomString(255) + "@gmail.com",
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)
	if err == nil {
		t.Errorf("Signup allows too long email address!")
	} else {
		t.Logf("Received error %v", err)
	}

}

func TestSignupNoPassword(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    "some@email.com",
		Password: "  ",
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)
	if err == nil {
		t.Errorf("Signup allows empty password !")
	} else {
		t.Logf("Received error %v", err)
	}

}

func TestSignupTooLongPassword(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    faker.Internet().Email(),
		Password: faker.RandomString(255),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)
	if err == nil {
		t.Errorf("Signup allows too big password !")
	} else {
		t.Logf("Received error %v", err)
	}

}

func TestSigningupMultipleTimes(t *testing.T) {

	email := faker.Internet().Email()

	params := usersvc.SignupParams{
		Email:    email,
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)
	if err != nil {
		t.Errorf("Unexpected error signing up %v", err)
		return
	}

	_, err = userSvcObj.SignupNewUser(ctx, &params)
	if err == nil {
		t.Errorf("Signup allows duplicate email address")
	} else {
		t.Logf("Received error %v", err)
	}
}

type SignupInterested struct {
	received_user_id string
	received_count   int64
}

func (s *SignupInterested) NewUser(ctx context.Context, u *db.User) error {
	fmt.Printf("Received new user notification ")
	s.received_user_id = u.ID
	s.received_count += 1
	return nil
}

func TestSignupExternalUser(t *testing.T) {

	var userInfo external.UserProvidedData
	userInfo.Email = faker.Internet().Email()
	userInfo.Name = faker.Name().Name()
	userInfo.Provider = "google"
	userInfo.Verified = true

	userSvcObj := usersvc.New(svc)

	token, err := userSvcObj.AuthenticateExternalUser(context.Background(),
		&userInfo)

	if err != nil {
		t.Errorf("Error signing up external user %v", err)
		return
	}
	t.Logf("Authenticate external user success token %v", utils.ToJSONString(token))

	ptoken, err := svc.GetJWTUtil().ParseJWTClaims(token.Token)
	if err != nil {
		t.Errorf("Error parsing token %v", err)
		return
	}
	t.Logf("parsed token %v", utils.ToJSONString(ptoken))

	user_id := ptoken.Subject
	userRec, err := userSvcObj.GetUser(context.Background(), user_id)
	if err != nil {
		t.Errorf("Error getting User Rec back %v", err)
		return
	}
	t.Logf("UserRec %v", utils.ToJSONString(userRec))

	if userRec.FirstName != userInfo.Name {
		t.Errorf("The user names does not match! %s %s ", userRec.FirstName, userInfo.Name)
	}

	if len(userRec.EncryptedPassword) <= 0 {
		t.Errorf("External User login password is not inittialized")
	}

	if !userRec.ConfirmedAt.Valid {
		t.Errorf("The external signed up user is not confirmed! ")
	}
}

func TestEmailConfirmation(t *testing.T) {

	email := faker.Internet().Email()

	params := usersvc.SignupParams{
		Email:    email,
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	user, err := userSvcObj.SignupNewUser(ctx, &params)
	if err != nil {
		t.Errorf("Unexpected error signing up %v", err)
		return
	}

	<-time.After(1 * time.Second)

	user2, err := userSvcObj.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("Error getting back the user rec %v", err)
		return
	}

	t.Logf("Users Email Confirmation Token %s ", user2.ConfirmationToken)
	if len(user2.ConfirmationToken) <= 0 {
		t.Error("Email Confirmation token is not sent! ")
	}
	if !user2.ConfirmationSentAt.Valid {
		t.Error("Email Confirmation sent time is not recorded ")
	}
	sentTime := user2.ConfirmationSentAt.Time

	err = userSvcObj.SendEmailConfirmationRequest(ctx, user.ID)
	if err != nil {
		t.Errorf("Sending confirmation request email %v", err)
		return
	}
	<-time.After(1 * time.Second)

	ur2_1, _ := userSvcObj.GetUser(ctx, user.ID)
	if !ur2_1.ConfirmationSentAt.Time.Equal(sentTime) {
		t.Error("Send email confirmation request is not checking the time difference ")
	}

	confirmCode := user2.ConfirmationToken

	resp, err := userSvcObj.ConfirmUserEmail(ctx, confirmCode)

	if err != nil {
		t.Errorf("Error confirming user %v", err)
	}

	t.Logf("ConfirmUserEmail resp %v", utils.ToJSONString(resp))

	ur3, err := userSvcObj.GetUser(ctx, user.ID)

	if !ur3.ConfirmedAt.Valid {
		t.Error("User email confirmation is not recorded ")
	}
}

func TestEmailConfirmationBadToken(t *testing.T) {
	ctx := context.Background()
	userSvcObj := usersvc.New(svc)

	confirmCode := faker.RandomString(8)
	_, err := userSvcObj.ConfirmUserEmail(ctx, confirmCode)
	assert.NotNil(t, err)

}

func TestPasswordReset(t *testing.T) {

	email := faker.Internet().Email()

	params := usersvc.SignupParams{
		Email:    email,
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	user, err := userSvcObj.SignupNewUser(ctx, &params)
	if err != nil {
		t.Errorf("Unexpected error signing up %v", err)
		return
	}

	user2, err := userSvcObj.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("Error getting back the user rec %v", err)
		return
	}

	if len(user2.RecoveryToken) != 0 {
		t.Error("Recovery token is not empty at signup")
	}

	if user2.RecoverySentAt.Valid {
		t.Error("Recovery token sent time is valid at signup!")
	}

	err = userSvcObj.InitResetPasswordRequest(ctx, user.Email)
	if err != nil {
		t.Errorf("Error in reset password request %v", err)
		return
	}
	<-time.After(1 * time.Second)

	ur3, err := userSvcObj.GetUser(ctx, user.ID)

	if len(ur3.RecoveryToken) <= 0 {
		t.Error("Recovery token is empty after reset password")
	}

	t.Logf("Reset password token %s", ur3.RecoveryToken)

	if !ur3.RecoverySentAt.Valid {
		t.Error("Recovery token sent time is invalid after reset password.")
	}

	err = userSvcObj.InitResetPasswordRequest(ctx, user.Email)
	if err == nil {
		t.Error("InitResetPasswordRequest allows resetting password in quick succession")
	}

	t.Logf("InitResetPasswordRequest rightfully errors when reset password is requested immediately %v", err)

	newpasswd := faker.RandomString(12)
	err = userSvcObj.ResetPassword(ctx, ur3.RecoveryToken, newpasswd)
	if err != nil {
		t.Errorf("Error in ResetPassword %v", err)
		return
	}
	atoken, err := userSvcObj.PasswordLogin(ctx,
		&usersvc.PasswordLoginParams{params.Email, params.Password})
	if err == nil {
		t.Error("Can login using old password")
	} else {
		t.Logf("Rightfully errors on old password err= %v", err)
	}
	if atoken != nil {
		t.Error("Can get access token with old password!")
	}

	atoken2, err := userSvcObj.PasswordLogin(ctx,
		&usersvc.PasswordLoginParams{email, newpasswd})
	if err != nil {
		t.Errorf("Error logging in with new password %v", err)
		return
	}
	if atoken2 == nil {
		t.Error("Error logging in with new password. Access token is null ")
	}
	t.Logf("Can login using new password token %v", atoken2)
}

func TestUpdateUserProfile(t *testing.T) {

	email := faker.Internet().Email()

	params := usersvc.SignupParams{
		Email:    email,
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	user, err := userSvcObj.SignupNewUser(ctx, &params)
	if err != nil {
		t.Errorf("Unexpected error signing up %v", err)
		return
	}

	fname := faker.Name().FirstName()
	err = userSvcObj.UpdateProfileField(ctx, user, "first_name", fname)
	if err != nil {
		t.Errorf("Unexpected error updating user profile %v", err)
		return
	}

	ur2, err := userSvcObj.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("Error getting back the user rec %v", err)
		return
	}

	if ur2.FirstName != fname {
		t.Errorf("User profile update - firstname does not match %s %s",
			ur2.FirstName, fname)
		return
	}

	lname := faker.Name().LastName()
	err = userSvcObj.UpdateProfileField(ctx, user, "last_name", lname)
	if err != nil {
		t.Errorf("Unexpected error updating user profile %v", err)
		return
	}

	ur3, err := userSvcObj.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("Error getting back the user rec %v", err)
		return
	}

	if ur3.LastName != lname {
		t.Errorf("User profile update - last name does not match %s %s",
			ur3.LastName, lname)
		return
	}

	newEmail := faker.Internet().Email()
	err = userSvcObj.UpdateProfileField(ctx, user, "email", newEmail)
	if err != nil {
		t.Errorf("Unexpected error updating user profile %v", err)
		return
	}

	<-time.After(1 * time.Second)

	ur4, err := userSvcObj.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("Error getting back the user rec %v", err)
		return
	}

	if ur4.EmailChange != newEmail {
		t.Errorf("User profile update - email change does not match %s %s",
			ur4.EmailChange, newEmail)
		return
	}

	if !ur4.EmailChangeSentAt.Valid {
		t.Error("User Profile update - email change sent is null")
		return
	}

	t.Logf("Email change token %v, new email %v ", ur4.EmailChangeToken, ur4.EmailChange)

	token := ur4.EmailChangeToken

	err = userSvcObj.CompleteEmailUpdate(ctx, token)
	if err != nil {
		t.Errorf("Unexpected error completing email update %v", err)
		return
	}

	ur5, err := userSvcObj.GetUser(ctx, user.ID)
	if err != nil {
		t.Errorf("Error getting back the user rec %v", err)
		return
	}

	if ur5.Email != newEmail {
		t.Errorf("Update profile didn't update email address %s | %s ", ur5.Email, newEmail)
	}

	if len(ur5.EmailChange) > 0 {
		t.Errorf("EmailChange field is not reset after email change done")
	}

	if ur5.EmailChangeSentAt.Valid {
		t.Errorf("EmailChangeSent date field is not reset after email change done")
	}

}

func TestUpdatePassword(t *testing.T) {
	email := faker.Internet().Email()

	params := usersvc.SignupParams{
		Email:    email,
		Password: faker.RandomString(8),
	}

	userSvcObj := usersvc.New(svc)

	ctx := context.Background()
	user, err := userSvcObj.SignupNewUser(ctx, &params)
	if err != nil {
		t.Errorf("Unexpected error signing up %v", err)
		return
	}
	atoken, err := userSvcObj.PasswordLogin(ctx,
		&usersvc.PasswordLoginParams{params.Email, params.Password})

	if err != nil {
		t.Errorf("Unexpected error logging in %v", err)
		return
	}
	t.Logf("Logged in - first time token %v ", atoken)

	new_password := faker.RandomString(10)
	err = userSvcObj.UpdatePassword(ctx, user, params.Password, new_password)
	if err != nil {
		t.Errorf("Unexpected error changing password %v", err)
		return
	}
	_, err = userSvcObj.PasswordLogin(ctx,
		&usersvc.PasswordLoginParams{params.Email, params.Password})
	if err == nil {
		t.Error("Can login with old password after chaning password")
		return
	}
	t.Logf("rightfully can't login with old password after changing password %v ", err)

	atoken2, err := userSvcObj.PasswordLogin(ctx,
		&usersvc.PasswordLoginParams{params.Email, new_password})
	if err != nil {
		t.Errorf("Error logging in with new password %v", err)
		return
	}
	t.Logf("Logged in with new password Token %v", atoken2)

}
