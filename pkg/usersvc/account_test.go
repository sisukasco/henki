package usersvc_test

import (
	"context"
	"encoding/json"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func createAppUser(email string) *usersvc.AppUser {
	au := &usersvc.AppUser{
		UserType: "pro",
		Reseller: "fastspring",
		AccountInfo: usersvc.AccountInfo{
			ID: faker.RandomString(12),
			Contact: usersvc.AccountContact{
				FirstName: faker.Name().FirstName(),
				LastName:  faker.Name().LastName(),
				Email:     email,
			},
			Language: "en",
			Country:  "USA",
		},
	}
	return au
}

func TestAccountUpdate(t *testing.T) {
	userSvcObj := usersvc.New(svc)
	ctx := context.Background()
	email := faker.Internet().SafeEmail()
	origPasswd := faker.RandomString(12)
	usr, err := userSvcObj.SignupNewUser(ctx, &usersvc.SignupParams{
		Email:    email,
		Password: origPasswd,
	})
	assert.Nil(t, err)

	au := createAppUser(email)

	resp, err := userSvcObj.UpdateUserAccount(ctx, email, au)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Greater(t, len(resp.UserID), 4)
	assert.Equal(t, usr.ID, resp.UserID)

	<-time.After(1 * time.Second)
	resetPwd, err := svc.DB.Q.GetResetPasswordOnConfirmation(ctx, resp.UserID)
	assert.Nil(t, err)

	assert.False(t, resetPwd)

	user, err := userSvcObj.GetUser(ctx, resp.UserID)
	assert.Greater(t, len(user.ConfirmationToken), 2)

	conf, err := userSvcObj.ConfirmUserEmail(ctx, user.ConfirmationToken)
	assert.Nil(t, err)

	t.Logf("confirm user email resp %v", utils.ToJSONString(conf))
	assert.Equal(t, "", conf.ResetPasswordToken)

	access, err := userSvcObj.PasswordLogin(ctx, &usersvc.PasswordLoginParams{
		Username: email,
		Password: origPasswd,
	})
	assert.Nil(t, err)

	t.Logf("User access %v", utils.ToJSONString(access))
}

func TestAccountUpdateUserDoesNotExist(t *testing.T) {
	userSvcObj := usersvc.New(svc)
	ctx := context.Background()
	email := faker.Internet().SafeEmail()

	au := createAppUser(email)

	resp, err := userSvcObj.UpdateUserAccount(ctx, email, au)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Greater(t, len(resp.UserID), 4)
	t.Logf("New user created ID %s ", resp.UserID)

	<-time.After(1 * time.Second)

	user, err := userSvcObj.GetUser(ctx, resp.UserID)
	assert.Greater(t, len(user.ConfirmationToken), 2)

	conf, err := userSvcObj.ConfirmUserEmail(ctx, user.ConfirmationToken)
	assert.Nil(t, err)

	t.Logf("confirm user email resp %v", utils.ToJSONString(conf))

	assert.Greater(t, len(conf.ResetPasswordToken), 5)

	newPassword := faker.RandomString(22)
	err = userSvcObj.ResetPassword(ctx, conf.ResetPasswordToken, newPassword)
	assert.Nil(t, err)

	access, err := userSvcObj.PasswordLogin(ctx, &usersvc.PasswordLoginParams{
		Username: email,
		Password: newPassword,
	})
	assert.Nil(t, err)

	t.Logf("User access %v", utils.ToJSONString(access))

	resetPwd, err := svc.DB.Q.GetResetPasswordOnConfirmation(ctx, resp.UserID)
	assert.Nil(t, err)

	assert.False(t, resetPwd)
}

/**
This is to generate a user account that can be used to test using the UI
*/
func TestGenerateAccountUpdateUserDoesNotExist(t *testing.T) {
	userSvcObj := usersvc.New(svc)
	ctx := context.Background()
	email := faker.Internet().SafeEmail()

	au := createAppUser(email)

	resp, err := userSvcObj.UpdateUserAccount(ctx, email, au)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Greater(t, len(resp.UserID), 4)
	t.Logf("New user created ID %s Email %s", resp.UserID, email)

}

func getUserType(ctx context.Context, userID string) (string, error) {
	user, err := svc.DB.Q.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}
	ui := &usersvc.UserInfo{}
	err = json.Unmarshal(user.UserInfo, ui)
	if err != nil {
		return "", err
	}

	//log.Printf("User Info\n%s\n", utils.ToJSONString(ui))
	return ui.AppUser.UserType, nil

}
func TestSwitchFromPaidToFree(t *testing.T) {
	userSvcObj := usersvc.New(svc)
	ctx := context.Background()
	email := faker.Internet().SafeEmail()
	origPasswd := faker.RandomString(12)
	usr, err := userSvcObj.SignupNewUser(ctx, &usersvc.SignupParams{
		Email:    email,
		Password: origPasswd,
	})
	assert.Nil(t, err)

	au := createAppUser(email)

	resp, err := userSvcObj.UpdateUserAccount(ctx, email, au)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Greater(t, len(resp.UserID), 4)
	assert.Equal(t, usr.ID, resp.UserID)

	utype, err := getUserType(ctx, usr.ID)
	assert.Nil(t, err)

	assert.Greater(t, len(utype), 1)

	remres, err := userSvcObj.SwitchToFreeAccount(ctx, email)
	assert.Nil(t, err)

	assert.Equal(t, remres.UserID, usr.ID)

	utype2, err := getUserType(ctx, usr.ID)
	assert.Nil(t, err)

	assert.Equal(t, len(utype2), 0)

}
