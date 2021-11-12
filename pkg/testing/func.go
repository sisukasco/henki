package testing

import (
	"context"
	"github.com/sisukasco/henki/pkg/db"
	"github.com/sisukasco/henki/pkg/service"
	"github.com/sisukasco/henki/pkg/usersvc"
	"syreclabs.com/go/faker"
	"testing"
)

func SignupUser(t *testing.T, svcParam *service.Service) (*db.User, error) {
	userSvcObj := usersvc.New(svcParam)

	params := usersvc.SignupParams{
		Email:    faker.Internet().Email(),
		Password: faker.RandomString(8),
	}

	ctx := context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)
	if err != nil {
		t.Errorf("Unexpected error signing up %v", err)
		return nil, err
	}

	return usersvc.GetUserByEmail(ctx, svcParam.DB.Q, params.Email)
}
