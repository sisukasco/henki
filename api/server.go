package api

import (
	"context"
	"fmt"
	"github.com/sisukasco/henki/pkg/service"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/knadh/koanf"
)

// Server provides an http.Server.
type Server struct {
	*http.Server
	svc    *service.Service
	webAPI *WebAPI
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer(konf *koanf.Koanf) (*Server, error) {
	log.Println("configuring server...")

	svc, err := service.NewService(konf)
	if err != nil {
		fmt.Printf("Error starting service %+v", err)
		os.Exit(1)
	}

	webAPI := NewWebAPI(svc)

	var addr string
	addr = ":" + konf.String("api.port")

	srv := http.Server{
		Addr:    addr,
		Handler: webAPI.Handler,
	}

	return &Server{&srv, svc, webAPI}, nil
}

func (srv *Server) Init() error {
	err := srv.svc.InitServer()
	if err != nil {
		return err
	}
	return nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	log.Println("starting server...")
	err := srv.Init()
	if err != nil {
		fmt.Printf("Error initializing server %+v", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)

	// teardown logic...
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	srv.webAPI.Shutdown()
	log.Println("Server gracefully stopped")
}
