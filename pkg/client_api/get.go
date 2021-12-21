package client_api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sisukasco/commons/crypto"
	"github.com/sisukasco/commons/http_utils"
)

type GetUserRequest struct {
	Client string `json:"client"`
	UserID string `json:"user_id"`
	APIKey string `json:"api_key"`
}

type GetUserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (ca *ClientApi) GetUser(w http.ResponseWriter, r *http.Request) {
	http_utils.HandleCall(ca.handleGetUser, w, r)
}

func (ca *ClientApi) VerifyClientRequest(clientName string, r *http.Request, message interface{}) error {

	signx, err := http_utils.ExtractSignature(r)
	if err != nil {
		return http_utils.ForbiddenError("Signature is required for this request")
	}

	clientKey := "apiClients." + clientName

	if !ca.svc.Konf.Exists(clientKey) {
		return http_utils.ForbiddenError("Client name is not recognised " + clientName)
	}

	verificationKey := ca.svc.Konf.String(clientKey + ".verificationKey")
	if len(verificationKey) < 5 {
		return errors.New("Client Signature is not set. ")
	}

	err = crypto.VerifySignature(message, verificationKey, signx)
	if err != nil {
		return http_utils.ForbiddenError("Bad signature for this request")
	}

	return nil
}

func (ca *ClientApi) handleGetUser(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	ur := &GetUserRequest{}

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
	user, err := ca.svc.DB.Q.GetUser(r.Context(), ur.UserID)
	if err != nil {
		return nil, err
	}

	resp := &GetUserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return resp, nil
}
