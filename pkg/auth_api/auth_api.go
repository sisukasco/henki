package auth_api

import (
	"github.com/sisukasco/henki/pkg/service"
	"github.com/sisukasco/henki/pkg/usersvc"
)

type AuthApi struct {
	usvc *usersvc.UserService
	svc  *service.Service
}

func NewAuthApi(svc *service.Service) *AuthApi {
	us := usersvc.New(svc)

	return &AuthApi{us, svc}
}

func (a *AuthApi) Shutdown() {
	a.usvc.Shutdown()
}
