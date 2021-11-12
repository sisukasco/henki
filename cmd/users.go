package main

import (
	"context"
	"fmt"
	"github.com/sisukasco/henki/pkg/db"
	"github.com/sisukasco/henki/pkg/service"

	"github.com/knadh/koanf"
)

func showUsers(konf *koanf.Koanf) {
	svc, err := service.NewService(konf)
	if err != nil {
		fmt.Printf("Error getting service %+v", err)
		return
	}
	users, err := svc.DB.Q.GetUsers(context.Background(), db.GetUsersParams{
		Limit:  20,
		Offset: 0,
	})
	if err != nil {
		fmt.Printf("Error getting users %+v", err)
		return
	}
	for _, u := range users {
		fmt.Printf("\n%s %s\n", u.Email, u.ID)
	}
}
