package usersvc

import (
	"context"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/commons/stringid"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"
	"github.com/sisukasco/henki/pkg/external"
	"strings"

	"github.com/prasanthmj/machine"
)

// SignupParams are the parameters the Signup endpoint accepts
type SignupParams struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Data     map[string]interface{} `json:"data"`
	Provider string                 `json:"-"`
}

func (this *UserService) SignupNewUser(ctx context.Context, params *SignupParams) (*db.User, error) {

	email := utils.CleanupString(params.Email)
	passwd := strings.TrimSpace(params.Password)

	if passwd == "" {
		return nil, http_utils.UnprocessableEntityError("Signup requires a valid password")
	}
	if len(passwd) > 250 {
		return nil, http_utils.UnprocessableEntityError("Too long password ")
	}
	if len(email) > 250 {
		return nil, http_utils.UnprocessableEntityError("Too long email address ")
	}
	err := utils.ValidateEmail(email)
	if err != nil {
		return nil, http_utils.UnprocessableEntityError("%v", err)
	}
	exists, err := DoesUserExist(ctx, this.svc.DB.Q, email)
	if err != nil {
		return nil, http_utils.InternalServerError("Couldn' Process Request").WithInternalError(err)
	}
	if exists {
		return nil, http_utils.BadRequestError("A user with this email address has already been registered")
	}
	user, err := NewUser(ctx, this.svc.DB.Q, email, passwd, params.Data)
	if err != nil {
		return nil, http_utils.InternalServerError("Couldn't create new user").WithInternalError(err)
	}

	job := machine.NewJob(&ConfirmationEmailTask{user.ID})
	this.svc.JQ.QueueUp(job)

	return user, nil
}

func (usvc *UserService) AuthenticateExternalUser(ctx context.Context,
	userInfo *external.UserProvidedData) (*AccessTokenResponse, error) {
	email := utils.CleanupString(userInfo.Email)

	exists, err := DoesUserExist(ctx, usvc.svc.DB.Q, email)
	if err != nil {
		return nil, http_utils.InternalServerError("Couldn't check for User").WithInternalError(err)
	}
	if exists {
		user, err := GetUserByEmail(ctx, usvc.svc.DB.Q, email)
		if err != nil {
			return nil, http_utils.InternalServerError("Couldn't get User Rec").WithInternalError(err)
		}
		return usvc.createAccessTokenForUser(ctx, user)
	} else {
		passwd := stringid.RandString(32)
		var params SignupParams

		params.Email = email
		params.Password = passwd
		params.Provider = userInfo.Provider

		user, err := usvc.SignupNewUser(ctx, &params)
		if err != nil {
			return nil, err
		}
		uparams := db.UpdateUserProfileParams{
			ID:        user.ID,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
			AvatarUrl: userInfo.AvatarURL,
		}
		if len(uparams.FirstName) <= 0 && len(userInfo.Name) > 0 {
			uparams.FirstName = userInfo.Name
		}

		err = usvc.svc.DB.Q.UpdateUserProfile(ctx, uparams)
		if err != nil {
			return nil, err
		}

		if userInfo.Verified {
			usvc.svc.DB.Q.ConfirmUserEmailByID(ctx, user.ID)
		}

		return usvc.createAccessTokenForUser(ctx, user)

	}

}
