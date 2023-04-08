package usersvc_test

import (
	"context"
	"testing"
	"time"

	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"
	"github.com/sisukasco/henki/pkg/usersvc"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestGettingUsersSignedUpNDaysAgoOneRec(t *testing.T) {
	userSvcObj := usersvc.New(svc)

	twoDaysBack := time.Now().AddDate(0, 0, -2)
	userRec := db.InsertUserCustomParams{
		ID:        faker.RandomString(8),
		Email:     faker.Internet().Email(),
		FirstName: faker.Name().FirstName(),
		LastName:  faker.Name().LastName(),
		CreatedAt: twoDaysBack,
		UpdatedAt: twoDaysBack,
	}

	ctx := context.Background()

	t.Logf("userRec %v", utils.ToJSONString(userRec))

	_, err := svc.DB.Q.InsertUserCustom(ctx, userRec)
	assert.Nil(t, err)

	users, err := userSvcObj.GetUsersCreatedNDaysAgo(ctx, 2)
	assert.Nil(t, err)

	t.Logf("Users %v", utils.ToJSONString(users))

	assert.Equal(t, len(users), 1)

	assert.Equal(t, users[0].Email, userRec.Email)

}

func insertUser(daysBack int32) (*db.InsertUserCustomParams, error) {
	ctx := context.Background()

	daysBackObj := time.Now().AddDate(0, 0, int(-1*daysBack))

	userRec := db.InsertUserCustomParams{
		ID:        faker.RandomString(8),
		Email:     faker.Internet().Email(),
		FirstName: faker.Name().FirstName(),
		LastName:  faker.Name().LastName(),
		CreatedAt: daysBackObj,
		UpdatedAt: daysBackObj,
	}
	_, err := svc.DB.Q.InsertUserCustom(ctx, userRec)
	if err != nil {
		return nil, err
	}
	return &userRec, nil
}

/*
*TODO: Fix this test later

func TestGettingUsersSignedUpNDaysAgoMultipleRec(t *testing.T) {
	userSvcObj := usersvc.New(svc)

	ctx := context.Background()

	userRec, err := insertUser(2)
	assert.Nil(t, err)

	_, err = insertUser(3)
	assert.Nil(t, err)

	_, err = insertUser(1)
	assert.Nil(t, err)

	users, err := userSvcObj.GetUsersCreatedNDaysAgo(ctx, 2)
	assert.Nil(t, err)

	t.Logf("Users %v", utils.ToJSONString(users))

	assert.Equal(t, len(users), 1)

	assert.Equal(t, users[0].Email, userRec.Email)

}
*/
