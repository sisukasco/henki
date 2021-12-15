package service

import (
	"log"

	"github.com/sisukasco/commons/auth"

	"github.com/knadh/koanf"
	"github.com/pkg/errors"
)

type Service struct {
	Konf  *koanf.Koanf
	DB    *DBConnection
	AuthM *auth.AuthMiddleware
}

func NewService(Konf *koanf.Koanf) (*Service, error) {

	log.Printf("Service starting... ")
	log.Printf("Connecting to Database ...")
	db, err := NewConnection(Konf.String("db.url"))
	if err != nil {
		return nil, errors.Wrap(err, "Service Initializing DB")
	}

	am := auth.NewAuthMiddleware(getJWTConf(Konf))

	return &Service{Konf: Konf, DB: db, AuthM: am}, nil
}

func (svc *Service) InitServer() error {

	return nil
}

func (this *Service) Close() {
	this.DB.Close()
}

func getJWTConf(konf *koanf.Koanf) *auth.JWTConfig {
	jconf := &auth.JWTConfig{
		Secret:        konf.String("jwt.secret"),
		ExpirySeconds: int32(konf.Int("jwt.expiry")),
		Aud:           konf.String("jwt.aud"),
	}

	return jconf
}
func (svc *Service) GetJWTUtil() *auth.JWTUtil {
	return auth.NewJWTUtil(getJWTConf(svc.Konf))
}
