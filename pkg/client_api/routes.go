package client_api

import (
	chi "github.com/go-chi/chi/v5"
)

func (ca *ClientApi) Routes(r *chi.Mux) {
	r.Post("/client/create/user", ca.CreateUser)
	r.Post("/client/get/user", ca.GetUser)
	r.Post("/client/get/users/ndays", ca.GetUsersCreatedNDaysAgo)
}
