package usersvc

import (
	"context"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(ctx context.Context, Q *db.Queries,
	email string, password string, userData utils.JSONMap) (*db.User, error) {

	email = utils.CleanupString(email)

	pw, err := hashPassword(password)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't create password hash")
	}
	user, err := Q.NewUser(ctx, db.NewUserParams{Email: email, EncryptedPassword: pw})
	if err != nil {
		return nil, errors.Wrap(err, "NewUser Error creating DB Entry")
	}
	return &user, nil
}
func (usvc *UserService) GetUser(ctx context.Context, strID string) (*db.User, error) {

	ID := strID

	user, err := usvc.svc.DB.Q.GetUser(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting User Record")
	}
	return &user, nil
}

func DoesUserExist(ctx context.Context, Q *db.Queries, email string) (bool, error) {

	email = utils.CleanupString(email)
	return Q.DoesUserExist(ctx, email)
}

func GetUserByEmail(ctx context.Context, Q *db.Queries, email string) (*db.User, error) {

	email = utils.CleanupString(email)

	user, err := Q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrapf(err, "GetUserByEmail %v", email)
	}
	return &user, nil
}

// hashPassword generates a hashed password from a plaintext string
func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}
