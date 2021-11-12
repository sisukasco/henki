package usersvc

import (
	"context"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/db"

	jwt "github.com/dgrijalva/jwt-go"
)

//TODO: Remove unnecessary data from JWT token (such as Email)
// Holds JWT claims
type MyJWTClaims struct {
	jwt.StandardClaims
	Email        string                 `json:"email"`
	AppMetaData  map[string]interface{} `json:"app_metadata"`
	UserMetaData map[string]interface{} `json:"user_metadata"`
}

// AccessTokenResponse represents an OAuth2 success response
type AccessTokenResponse struct {
	Token        string `json:"access_token"`
	TokenType    string `json:"token_type"` // Bearer
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
type PasswordLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (usvc *UserService) RenewRefreshToken(ctx context.Context,
	current_token string) (*AccessTokenResponse, error) {

	if current_token == "" {
		return nil, http_utils.OauthError("invalid_request", "refresh_token required")
	}

	user, err := usvc.svc.DB.Q.GetUserFromRefreshToken(ctx, current_token)
	if err != nil {
		return nil, http_utils.OauthError("invalid_grant", "Invalid Refresh Token").WithInternalMessage("Possible abuse attempt: %v", err)
	}
	usvc.svc.DB.Q.ClearRefreshTokens(ctx, user.ID)

	return usvc.createAccessTokenForUser(ctx, &user)
}

func (usvc *UserService) PasswordLogin(ctx context.Context, params *PasswordLoginParams) (*AccessTokenResponse, error) {

	email := utils.CleanupString(params.Username)

	user, err := GetUserByEmail(ctx, usvc.svc.DB.Q, email)
	if err != nil {
		return nil, http_utils.OauthError("access_denied", "no such user")
	}
	if !user.Authenticate(params.Password) {
		return nil, http_utils.OauthError("invalid_grant", "Invalid Password")
	}

	return usvc.createAccessTokenForUser(ctx, user)
}

func (usvc *UserService) AdminLoginAs(ctx context.Context, email string) (*AccessTokenResponse, error) {
	email = utils.CleanupString(email)

	user, err := GetUserByEmail(ctx, usvc.svc.DB.Q, email)
	if err != nil {
		return nil, http_utils.OauthError("access_denied", "no such user")
	}
	return usvc.createAccessTokenForUser(ctx, user)
}

func (usvc *UserService) createAccessTokenForUser(ctx context.Context,
	user *db.User) (*AccessTokenResponse, error) {

	refreshToken, err := createRefreshToken(ctx, usvc.svc.DB.Q, user.ID)
	if err != nil {
		return nil, http_utils.InternalServerError("Error creating token for user").WithInternalError(err)
	}

	accessToken, err := usvc.svc.GetJWTUtil().GenerateAccessToken(user.ID)
	if err != nil {
		return nil, http_utils.InternalServerError("Error creating access token for user").WithInternalError(err)
	}

	expires := usvc.svc.Konf.Int("jwt.expiry")

	if expires < 200 {
		expires = 200
	}

	return &AccessTokenResponse{
		Token:        accessToken,
		TokenType:    "bearer",
		ExpiresIn:    expires,
		RefreshToken: refreshToken,
	}, nil
}
