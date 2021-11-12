package auth_api

import (
	"encoding/json"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/auth_utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"log"
	"net/http"
)

type UpdateProfileField struct {
	FieldName string `json:"name"`
	Value     string `json:"value"`
}

func (a *AuthApi) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, err := auth_utils.GetUserFromRequest(a.svc, r)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	params := &UpdateProfileField{}

	jsonDecoder := json.NewDecoder(r.Body)

	err = jsonDecoder.Decode(params)

	if err != nil {
		err = http_utils.InternalServerError("Error decoding request").WithInternalError(err)
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	err = a.usvc.UpdateProfileField(r.Context(), user, params.FieldName, params.Value)

	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}
	err = http_utils.SendJSON(w, http.StatusOK, &http_utils.StatusResponse{"ok"})
	if err != nil {
		log.Printf("Error sending response %v", err)
	}

}

type ServiceEndPoints map[string]string

type UserResponse struct {
	Status         string           `json:"status"`
	ID             string           `json:"id"`
	Email          string           `json:"email"`
	EmailConfirmed bool             `json:"email_confirmed"`
	FirstName      string           `json:"first_name"`
	LastName       string           `json:"last_name"`
	AvatarURL      string           `json:"avatar_url"`
	EndPoints      ServiceEndPoints `json:"endpoints"`
	PaidUser       bool             `json:"paid_user"`
}

func (a *AuthApi) GetUser(w http.ResponseWriter, r *http.Request) {

	user, err := auth_utils.GetUserFromRequest(a.svc, r)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
	}
	resp := &UserResponse{Status: `ok`,
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarUrl,
	}
	if user.ConfirmedAt.Valid {
		resp.EmailConfirmed = true
	}

	ui := &usersvc.UserInfo{}
	err = json.Unmarshal(user.UserInfo, ui)
	if err != nil {
		log.Printf("Error unmarshalling UserInfo %v", err)
	} else {
		log.Printf("UserType is %v", ui.AppUser.UserType)
		if len(ui.AppUser.UserType) > 0 {
			resp.PaidUser = true
		}
	}

	resp.EndPoints = a.svc.Konf.StringMap("services")

	err = http_utils.SendJSON(w, http.StatusOK, resp)

	if err != nil {
		log.Printf("Error sending GetUser response %v", err)
	}
}
