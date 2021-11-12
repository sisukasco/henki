package auth_utils

import (
	"github.com/pkg/errors"
	"github.com/sisukasco/commons/auth"
	"github.com/sisukasco/henki/pkg/db"
	"github.com/sisukasco/henki/pkg/service"
	"net/http"
	"regexp"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

func GetUserFromRequest(svc *service.Service, r *http.Request) (*db.User, error) {
	userID, err := auth.GetUserIDFromContext(r)
	if err != nil {
		return nil, errors.Wrap(err, "No User in the context")
	}

	user, err := svc.DB.Q.GetUser(r.Context(), userID)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting User Record")
	}
	return &user, nil
}
