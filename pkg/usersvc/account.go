package usersvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sisukasco/commons/stringid"
	"github.com/sisukasco/commons/utils"
)

type UserAccountResponse struct {
	UserID string `json:"user_id"`
}

func (usvc *UserService) UpdateUserAccount(ctx context.Context, email string, appUser *AppUser) (*UserAccountResponse, error) {

	email = utils.CleanupString(email)
	exists, err := usvc.svc.DB.Q.DoesUserExist(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return usvc.updateExistingUserAccount(ctx, email, appUser)
	} else {
		return usvc.createNewUserAccount(ctx, email, appUser)
	}
}

func (usvc *UserService) updateExistingUserAccount(ctx context.Context, email string, appUser *AppUser) (*UserAccountResponse, error) {
	user, err := usvc.svc.DB.Q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = usvc.UpdateAppUser(ctx, user.ID, appUser, false)
	if err != nil {
		return nil, err
	}

	return &UserAccountResponse{UserID: user.ID}, nil
}

func (usvc *UserService) createNewUserAccount(ctx context.Context, email string, appUser *AppUser) (*UserAccountResponse, error) {
	//Createuser with option to change password on first login
	params := SignupParams{
		Email:    email,
		Password: stringid.RandString(22),
	}
	user, err := usvc.SignupNewUser(ctx, &params)
	if err != nil {
		return nil, err
	}
	err = usvc.UpdateAppUser(ctx, user.ID, appUser, true)
	if err != nil {
		return nil, err
	}
	return &UserAccountResponse{UserID: user.ID}, nil
}

func (usvc *UserService) SwitchToFreeAccount(ctx context.Context, email string) (*UserAccountResponse, error) {

	email = utils.CleanupString(email)

	user, err := usvc.svc.DB.Q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SwitchToFreeAccount -> No user with email %s ", email))
	}
	if len(user.UserInfo) <= 5 {
		return &UserAccountResponse{UserID: user.ID}, nil
	}
	ui := &UserInfo{}
	err = json.Unmarshal(user.UserInfo, ui)
	if err != nil {
		return nil, err
	}
	ui.AppUser.UserType = ""
	ui.AppUser.Updated = time.Now()

	err = usvc.UpdateAppUser(ctx, user.ID, &ui.AppUser, false)
	if err != nil {
		return nil, err
	}

	return &UserAccountResponse{UserID: user.ID}, nil
}
