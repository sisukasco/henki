package client_api

import (
	"github.com/sisukasco/henki/pkg/service"
	"github.com/sisukasco/henki/pkg/usersvc"
)

type ClientApi struct {
	svc  *service.Service
	usvc *usersvc.UserService
}

func NewClientApi(svc *service.Service, usvc *usersvc.UserService) *ClientApi {

	return &ClientApi{svc, usvc}
}
