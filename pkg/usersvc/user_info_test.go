package usersvc_test

import (
	"context"
	"encoding/json"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestUpdateAppUser(t *testing.T) {
	userSvcObj := usersvc.New(svc)

	params := usersvc.SignupParams{
		Email:    faker.Internet().Email(),
		Password: faker.RandomString(12),
	}

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)
	assert.Nil(t, err)

	user, err := svc.DB.Q.GetUserByEmail(ctx, params.Email)
	assert.Nil(t, err)

	var au usersvc.AppUser
	au.UserType = "pro"
	au.Reseller = "fastspring"
	au.Updated = time.Now()
	au.AccountInfo.ID = faker.RandomString(12)
	au.AccountInfo.Country = "USA"
	au.AccountInfo.Contact.Email = faker.Internet().Email()
	au.AccountInfo.Contact.FirstName = faker.Name().FirstName()
	au.AccountInfo.Contact.LastName = faker.Name().LastName()

	err = userSvcObj.UpdateAppUser(ctx, user.ID, &au, true)
	assert.Nil(t, err)

	bui, err := svc.DB.Q.GetUserInfo(ctx, user.ID)
	assert.Nil(t, err)

	ui := &usersvc.UserInfo{}

	err = json.Unmarshal(bui, ui)
	assert.Nil(t, err)

	assert.Equal(t, ui.AppUser.AccountInfo.ID, au.AccountInfo.ID)
	assert.Equal(t, ui.AppUser.AccountInfo.Contact.Email, au.AccountInfo.Contact.Email)
	assert.True(t, ui.ResetPasswordOnConfirmation)

	uu, err := svc.DB.Q.GetUser(ctx, user.ID)
	log.Printf("User rec\n%v\n", utils.ToJSONString(uu))

	assert.Equal(t, uu.FirstName, au.AccountInfo.Contact.FirstName)
	assert.Equal(t, uu.LastName, au.AccountInfo.Contact.LastName)

}
