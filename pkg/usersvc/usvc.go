package usersvc

import (
	"log"
	"sync"

	"github.com/sisukasco/henki/pkg/service"
)

type UserService struct {
	svc *service.Service
	wg  sync.WaitGroup
}

func New(svc *service.Service) *UserService {
	usvc := &UserService{svc: svc}
	return usvc
}
func (usvc *UserService) Shutdown() {
	log.Printf("UserService shutting down ...")
	usvc.wg.Wait()
	log.Printf("UserService shutdown.")
}
