package api

import (
	"net/http"
	"time"

	"github.com/sisukasco/henki/pkg/auth_api"
	"github.com/sisukasco/henki/pkg/client_api"
	"github.com/sisukasco/henki/pkg/service"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/knadh/koanf"
)

const (
	audHeaderName   = "X-JWT-AUD"
	useCookieHeader = "x-use-cookie"
)

type WebAPI struct {
	authAPI   *auth_api.AuthApi
	pingAPI   *PingApi
	Handler   *chi.Mux
	clientAPI *client_api.ClientApi
}

func NewWebAPI(svc *service.Service) *WebAPI {

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	//r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15 * time.Second))

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(corsConfig(svc.Konf)))

	pingAPI := NewPingApi(svc)
	pingAPI.routes(r)

	authAPI := auth_api.NewAuthApi(svc)
	authAPI.Routes(r)

	clientAPI := client_api.NewClientApi(svc, authAPI.GetUserService())

	if svc.Konf.Bool("apiClients.enable") {
		clientAPI.Routes(r)
	}

	webAPI := &WebAPI{authAPI, pingAPI, r, clientAPI}

	return webAPI
}

func corsConfig(konf *koanf.Koanf) cors.Options {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	opp := cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: konf.Strings("client.origins"),
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", audHeaderName, useCookieHeader},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	}

	return opp
}

func (webapi *WebAPI) Shutdown() {
	webapi.authAPI.Shutdown()
}
