package usersvc_test

import (
	"context"
	"encoding/json"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"testing"

	"syreclabs.com/go/faker"
)

func TestSignupAndLogin(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    faker.Internet().Email(),
		Password: faker.RandomString(20),
	}

	userSvcObj := usersvc.New(svc)

	ctx = context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)

	if err != nil {
		t.Errorf("Error signingup new user %v ", err)
		return
	}

	loginParams := usersvc.PasswordLoginParams{
		Username: params.Email,
		Password: params.Password,
	}
	token, err := userSvcObj.PasswordLogin(ctx, &loginParams)
	if err != nil {
		t.Errorf("Error logging in new user %v ", err)
		return
	}
	t.Logf("TestSignupAndLogin success %v", token.Token)

	loginParams2 := usersvc.PasswordLoginParams{
		Username: params.Email,
		Password: "  ",
	}
	token, err = userSvcObj.PasswordLogin(ctx, &loginParams2)
	if err == nil {
		t.Errorf("user logging in using empty password ")
		return
	} else {
		t.Logf("attempt logging in using empty password %v ", err)
	}

	loginParams3 := usersvc.PasswordLoginParams{
		Username: " " + params.Email + "  ",
		Password: params.Password,
	}
	token, err = userSvcObj.PasswordLogin(ctx, &loginParams3)
	if err != nil {
		t.Errorf("Error logging in using email with extra spaces ")
		return
	}
}

func TestAccessTokenAccess(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    faker.Internet().Email(),
		Password: faker.RandomString(20),
	}

	userSvcObj := usersvc.New(svc)

	ctx = context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)

	if err != nil {
		t.Errorf("Error signingup new user %v ", err)
		return
	}

	loginParams := usersvc.PasswordLoginParams{
		Username: params.Email,
		Password: params.Password,
	}
	token, err := userSvcObj.PasswordLogin(ctx, &loginParams)
	if err != nil {
		t.Errorf("Error logging in new user %v ", err)
		return
	}

	strToken, _ := json.MarshalIndent(token, "", "  ")

	t.Logf("TestAccessTokenAccess success %v", string(strToken))

	accessToken, err := svc.GetJWTUtil().ParseJWTClaims(token.Token)
	if err != nil {
		t.Errorf("Parsing JWT claims %v ", err)
		return
	}

	strAccessToken, _ := json.MarshalIndent(accessToken, "", "  ")

	t.Logf("TestAccessToken Decoded %v", string(strAccessToken))
}

//Test With expired tokens
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9eyJhdWQiOiJhcGkuZG9ja2Zvcm0uY29tIiwiZXhwIjoxNTg5NTQxNTEzLCJzdWIiOiJ6Wlhna3c4ZzJNIiwiZW1haWwiOiJjYWxlQGxlYW5ub25hYmJvdHQubmFtZSIsImFwcF9tZXRhZGF0YSI6bnVsbCwidXNlcl9tZXRhZGF0YSI6bnVsbH0.ATc5mCC20mu5_usHyxFksL8JlfvORFBxRsf_FStGzdw

func TestMakingAccessToken(t *testing.T) {

	fake_token2 := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhcGkuZG9ja2Zvcm0uY29tIiwiZXhwIjoxNTg5NTQyNDI4LCJzdWIiOiJBNWtnWTVNZzRtIiwiZW1haWwiOiJtaW5uaWVAZ3VzaWtvd3NraS5vcmciLCJhcHBfbWV0YWRhdGEiOm51bGwsInVzZXJfbWV0YWRhdGEiOm51bGx9.DXyVjOzvKmg8cLEclwurFYaPlH_7JBt2_uKpv-UyONM`
	//^ expired token

	fake_token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9eyJhdWQiOiJhcGkuZG9ja2Zvcm0uY29tIiwiZXhwIjoxNTg5NTQxNTEzLCJzdWIiOiJ6Wlhna3c4ZzJNIiwiZW1haWwiOiJjYWxlQGxlYW5ub25hYmJvdHQubmFtZSIsImFwcF9tZXRhZGF0YSI6bnVsbCwidXNlcl9tZXRhZGF0YSI6bnVsbH0.ATc5mCC20mu5_usHyxFksL8JlfvORFBxRsf_FStGzdw`

	accessToken, err := svc.GetJWTUtil().ParseJWTClaims(fake_token)
	if err == nil {
		t.Errorf("Fake access token successful! %v ", accessToken)
		return
	}

	t.Logf("TestMakingAccessToken test success error %v", err)

	accessToken, err = svc.GetJWTUtil().ParseJWTClaims(fake_token2)
	if err == nil {
		t.Errorf("Expired access token successful! %v ", accessToken)
		return
	}
	t.Logf("Test using expired access token test success error %v", err)
}

func TestRefreshTokenRenewal(t *testing.T) {
	params := usersvc.SignupParams{
		Email:    faker.Internet().Email(),
		Password: faker.RandomString(20),
	}

	userSvcObj := usersvc.New(svc)

	ctx = context.Background()
	_, err := userSvcObj.SignupNewUser(ctx, &params)

	if err != nil {
		t.Errorf("Error signingup new user %v ", err)
		return
	}

	loginParams := usersvc.PasswordLoginParams{
		Username: params.Email,
		Password: params.Password,
	}
	token, err := userSvcObj.PasswordLogin(ctx, &loginParams)
	if err != nil {
		t.Errorf("Error logging in new user %v ", err)
		return
	}

	refresh_token := token.RefreshToken

	t.Logf("Trying to renew refresh token %s", refresh_token)

	token2, err := userSvcObj.RenewRefreshToken(ctx, refresh_token)
	if err != nil {
		t.Errorf("Error renewing the refresh token %v", err)
		return
	}

	t.Logf("New token is %v ", utils.ToJSONString(token2))

	_, err = userSvcObj.RenewRefreshToken(ctx, refresh_token)
	if err == nil {
		t.Errorf("Can renew an already renewed token! ")
		return
	}

	t.Logf("Success -> Cant renew an already renewed refresh token ")
}
