package usersvc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sisukasco/henki/pkg/db"
)

func (usvc *UserService) GetUsersCreatedNDaysAgo(ctx context.Context, days int32) ([]db.GetUsersSignedUpNDaysAgoRow, error) {
	sDays := sql.NullString{
		String: fmt.Sprintf("%d", days),
		Valid:  true,
	}
	users, err := usvc.svc.DB.Q.GetUsersSignedUpNDaysAgo(ctx, sDays)
	if err != nil {
		return nil, err
	}

	return users, nil
}
