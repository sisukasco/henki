package api

import (
	"context"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/service"
	"log"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
)

type PingApi struct {
	svc *service.Service
}

func NewPingApi(svc *service.Service) *PingApi {
	return &PingApi{svc}
}

func (a *PingApi) routes(r *chi.Mux) {
	r.Get("/ping", a.pingTest)
}

func (a *PingApi) pingTest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Ping command received ...")
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	err := a.svc.DB.Ping(ctx)
	if err != nil {
		err = http_utils.InternalServerError("DB Ping failed").WithInternalError(err)
		http_utils.SendErrorResponse(err, w, r)
	}
	w.Write([]byte("pong"))
}
