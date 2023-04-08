package client_api

import (
	"encoding/json"
	"net/http"

	"github.com/sisukasco/commons/http_utils"
)

type GetUsersNDaysRequest struct {
	Client string `json:"client"`
	APIKey string `json:"api_key"`
	Days   int32  `json:"days"`
}

func (ca *ClientApi) GetUsersCreatedNDaysAgo(w http.ResponseWriter, r *http.Request) {
	http_utils.HandleCall(ca.handleGetUsersCreatedNDaysAgo, w, r)
}

func (ca *ClientApi) handleGetUsersCreatedNDaysAgo(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	ur := &GetUsersNDaysRequest{}

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

	users, err := ca.usvc.GetUsersCreatedNDaysAgo(r.Context(), ur.Days)
	if err != nil {
		return nil, err
	}

	return users, nil
}
