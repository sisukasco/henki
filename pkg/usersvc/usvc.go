package usersvc

import (
	"github.com/sisukasco/henki/pkg/service"
)

type UserService struct {
	svc *service.Service
}

func New(svc *service.Service) *UserService {
	usvc := &UserService{svc: svc}
	usvc.registerTaskTypes()
	return usvc
}
