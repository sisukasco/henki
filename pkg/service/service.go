package service

import (
	"github.com/sisukasco/commons/auth"
	"github.com/sisukasco/commons/redis"
	"log"

	"github.com/knadh/koanf"
	"github.com/pkg/errors"
	"github.com/prasanthmj/machine"
)

type Service struct {
	Konf  *koanf.Koanf
	DB    *DBConnection
	JQ    *machine.JobQueue
	Redis *redis.Redis
	AuthM *auth.AuthMiddleware
}

func NewService(Konf *koanf.Koanf) (*Service, error) {

	log.Printf("Service starting... ")
	log.Printf("Connecting to Database ...")
	db, err := NewConnection(Konf.String("db.url"))
	if err != nil {
		return nil, errors.Wrap(err, "Service Initializing DB")
	}

	machineRedisURL := Konf.String("machine.redis.url")

	jq, err := machine.New(machineRedisURL)
	if err != nil {
		return nil, errors.Wrap(err, "Service Initializing JQ")
	}

	redisURL := Konf.String("redis.url")
	red, err := redis.New(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "Service Initializing Redis")
	}

	am := auth.NewAuthMiddleware(getJWTConf(Konf))

	return &Service{Konf: Konf, DB: db, JQ: jq, Redis: red, AuthM: am}, nil
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
