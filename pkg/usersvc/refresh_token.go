package usersvc

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"
)

func createRefreshToken(ctx context.Context, Q *db.Queries, user_id string) (string, error) {
	rtoken := utils.SecureToken()

	err := Q.NewRefreshToken(ctx, db.NewRefreshTokenParams{
		UserID: user_id,
		Token:  rtoken,
	})
	if err != nil {
		return "", errors.Wrap(err, "creating refresh token")
	}
	return rtoken, nil
}
