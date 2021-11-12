package auth_api

import (
	chi "github.com/go-chi/chi/v5"
)

func (a *AuthApi) Routes(r *chi.Mux) {
	r.Post("/signup", a.Signup)
	r.Post("/token", a.GetToken)
	r.Post("/confirm", a.ConfirmEmail)
	r.Post("/reset/init", a.SendResetPassword)
	r.Post("/reset/update", a.ResetPassword)

	r.Route("/user", func(r chi.Router) {
		r.Use(a.svc.AuthM.RequireAuth)
		r.Get("/", a.GetUser)
		r.Post("/resend/conf/email", a.ResendConfirmEmail)
		r.Post("/profile", a.UpdateProfile)
		r.Post("/update/password", a.UpdatePassword)
	})

	r.Route("/api-key", func(r chi.Router) {
		r.Use(a.svc.AuthM.RequireAuth)

		r.Get("/", a.GetApiKeys)
		r.Post("/", a.CreateAPIKey)
		r.Delete("/", a.DeleteAPIKey)
	})
	r.Post("/email/update", a.CommitEmailUpdate)
	// External Login
	r.Get("/authorize", a.ExternalProviderRedirect)
	r.Get("/callback", a.ExternalProviderCallback)

	//change user account type
	r.Post("/account/update", a.UpdateAccount)
	r.Post("/account/remove", a.RemoveAccount)

}
