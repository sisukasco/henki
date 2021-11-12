package auth_api

import (
	"encoding/json"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/auth_utils"
	"github.com/sisukasco/henki/pkg/db"
	"log"
	"net/http"
)

func (a *AuthApi) GetApiKeys(w http.ResponseWriter, r *http.Request) {
	user, err := auth_utils.GetUserFromRequest(a.svc, r)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}
	keys, err := a.usvc.GetAPIKeys(r.Context(), user.ID)
	if err != nil {
		e := http_utils.InternalServerError("Error getting API Keys for user ").WithInternalError(err)
		http_utils.SendErrorResponse(e, w, r)
		return
	}
	if keys == nil {
		keys = make([]db.GetAPIKeysRow, 0)
	}
	err = http_utils.SendJSON(w, http.StatusOK, keys)

	if err != nil {
		log.Printf("Error sending GetApiKeys response %v", err)
	}
}

type CreateAPIKeyResponse struct {
	Key string `json:"key"`
}

func (a *AuthApi) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	user, err := auth_utils.GetUserFromRequest(a.svc, r)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	key, err := a.usvc.NewAPIKey(r.Context(), user.ID)

	if err != nil {
		e := http_utils.InternalServerError("Error creating API Keys for user ").WithInternalError(err)
		http_utils.SendErrorResponse(e, w, r)
		return
	}

	resp := CreateAPIKeyResponse{key}

	err = http_utils.SendJSON(w, http.StatusOK, resp)

	if err != nil {
		log.Printf("Error sending GetApiKeys response %v", err)
	}
}

type deleteAPIKeyParams struct {
	Key string `json:"key"`
}

func (a *AuthApi) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	user, err := auth_utils.GetUserFromRequest(a.svc, r)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	params := &deleteAPIKeyParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err = jsonDecoder.Decode(params)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	err = a.usvc.DeleteAPIKey(r.Context(), params.Key, user.ID)
	if err != nil {
		e := http_utils.InternalServerError("Error DeleteAPIKey API Keys for user ").WithInternalError(err)
		http_utils.SendErrorResponse(e, w, r)
		return
	}

	resp := http_utils.StatusResponse{"ok"}

	err = http_utils.SendJSON(w, http.StatusOK, resp)

	if err != nil {
		log.Printf("Error sending GetApiKeys response %v", err)
	}
}
