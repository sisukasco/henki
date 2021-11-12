package usersvc_test

import (
	"context"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"syreclabs.com/go/faker"
	"testing"
)

func TestNewUser(t *testing.T) {

	email := faker.Internet().Email()

	testOneUser(t, email, "a random password ")

	email = `UpperCase@Website.com`

	testOneUser(t, email, "a random password ")

	email = ` email-Spaces@Website.com   `

	testOneUser(t, email, "a random password ")

	for i := 0; i < 10; i++ {
		emailx := faker.Internet().Email()
		passwd := faker.Internet().Email()
		testOneUser(t, emailx, passwd)
	}
	t.Logf("Test %v completed ", t.Name())

}

func testOneUser(t *testing.T, email string, password string) bool {

	ctx := context.Background()
	user, err := usersvc.NewUser(ctx, svc.DB.Q, email, password,
		utils.JSONMap{})
	if err != nil {
		t.Errorf("Error creating user(%s) %#v", email, err)
		return false
	}
	t.Logf("Created new user %v", utils.ToJSONString(user))

	exists, err := usersvc.DoesUserExist(ctx, svc.DB.Q, email)
	if err != nil {
		t.Errorf("DoesUserExist throws error %s %v", email, err)
		return false
	}
	if !exists {
		t.Errorf("Error creating user. DoesUserExist returns false %v ", email)

	}

	userSvcObj := usersvc.New(svc)

	id := user.ID

	u1, err := userSvcObj.GetUser(ctx, id)
	emailupd := utils.CleanupString(email)

	if u1.Email != emailupd {
		t.Errorf("NewUser email does not match %v %v ", u1.Email, emailupd)
		return false
	}

	u2, err := usersvc.GetUserByEmail(ctx, svc.DB.Q, email)
	if err != nil {
		t.Errorf("NewUser getting by email didn't work %v ", email)
	}
	if u2.Email != emailupd {
		t.Errorf("NewUser email does not match %v %v ", u2.Email, emailupd)
	}
	return true
}
