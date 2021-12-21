package client_api

import (
	"encoding/json"
	"net/http"

	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/usersvc"
)

type CreateUserRequest struct {
	Client string `json:"client"`
	APIKey string `json:"api_key"`
	Email  string `json:"email"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func (ca *ClientApi) CreateUser(w http.ResponseWriter, r *http.Request) {
	http_utils.HandleCall(ca.handleCreateUser, w, r)
}
func (ca *ClientApi) handleCreateUser(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	ur := &CreateUserRequest{}

	jsonDecoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := jsonDecoder.Decode(ur)
	if err != nil {
		return nil, err
	}

	err = ca.VerifyClientRequest(ur.Client, r, ur)
	if err != nil {
		return nil, err
	}

	params := &usersvc.SignupParams{
		Email:    ur.Email,
		Password: utils.SecureToken(),
	}

	user, err := ca.usvc.SignupNewUser(r.Context(), params)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{UserID: user.ID}, nil
}
