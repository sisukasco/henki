package usersvc

import (
	"context"
	"encoding/json"
	"github.com/sisukasco/henki/pkg/db"
	"log"
	"time"
)

type AccountContact struct {
	FirstName string `json:"first"`
	LastName  string `json:"last"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Phone     string `json:"phone"`
}
type AccountInfo struct {
	ID       string         `json:"id"`
	Contact  AccountContact `json:"contact"`
	Language string         `json:"language"`
	Country  string         `json:"country"`
}

type AppUser struct {
	UserType    string      `json:"user_type"`
	Updated     time.Time   `json:"updated"`
	Reseller    string      `json:"reseller"`
	AccountInfo AccountInfo `json:"account_info"`
}

type UserInfo struct {
	AppUser                     AppUser `json:"app_user"`
	ResetPasswordOnConfirmation bool    `json:"reset_password_on_confirmation"`
}

func (usvc *UserService) updateResetPasswordOnConfirmation(ctx context.Context, userID string, reset bool) error {
	bRRUpdate, err := json.Marshal(reset)
	if err != nil {
		return err
	}
	err = usvc.svc.DB.Q.UpdateUserInfo(ctx,
		db.UpdateUserInfoParams{
			ID:          userID,
			Path:        "{reset_password_on_confirmation}",
			Replacement: bRRUpdate,
		})
	return err
}

func (usvc *UserService) UpdateAppUser(ctx context.Context, userID string,
	appUser *AppUser, resetPasswordOnConfirmation bool) error {
	bAppuser, err := json.Marshal(appUser)
	if err != nil {
		return err
	}
	err = usvc.svc.DB.Q.UpdateUserInfo(ctx,
		db.UpdateUserInfoParams{
			ID:          userID,
			Path:        "{app_user}",
			Replacement: bAppuser,
		})
	if err != nil {
		return err
	}

	if resetPasswordOnConfirmation {
		err = usvc.updateResetPasswordOnConfirmation(ctx, userID, true)
		if err != nil {
			log.Printf("Error updating resetPasswordOnConfirmation flag %v", err)
		}
	}

	uu, err := usvc.svc.DB.Q.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if len(uu.FirstName) <= 0 && len(uu.LastName) <= 0 {
		if len(appUser.AccountInfo.Contact.FirstName) > 0 ||
			len(appUser.AccountInfo.Contact.LastName) > 0 {
			err := usvc.svc.DB.Q.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
				ID:        userID,
				FirstName: appUser.AccountInfo.Contact.FirstName,
				LastName:  appUser.AccountInfo.Contact.LastName,
				AvatarUrl: uu.AvatarUrl,
			})
			if err != nil {
				log.Printf("Error updating the user profile %v", err)
			}
		}

	}
	return nil
}
