package db

import (
	"context"
	"github.com/sisukasco/commons/stringid"

	"golang.org/x/crypto/bcrypt"
)

func (u User) Authenticate(passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(passwd))
	return err == nil
}

const idLength = 10

func (q *Queries) generateUniqueUserID(ctx context.Context) string {
	id := stringid.RandID(idLength)
	t := 0
	for ; t < 100; t++ {
		exists, err := q.DoesUserIDExist(ctx, id)
		if exists == false && err == nil {
			break
		}
		len := (t / 10) + idLength
		id = stringid.RandID(len)
	}
	return id
}

type NewUserParams struct {
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`
}

func (q *Queries) NewUser(ctx context.Context,
	arg NewUserParams) (User, error) {
	userID := q.generateUniqueUserID(ctx)

	return q.createNewUser(ctx, createNewUserParams{userID,
		arg.Email, arg.EncryptedPassword})
}
