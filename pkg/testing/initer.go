package testing

import (
	"log"
	"os"
	"testing"

	"github.com/sisukasco/henki/pkg/service"

	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	kfile "github.com/knadh/koanf/providers/file"
)

func InitService(m *testing.M, before func(s *service.Service), after func()) {

	//setup
	var err error

	konf := koanf.New(".")
	konf.Load(kfile.Provider("../../conf-dev.yaml"), kyaml.Parser())
	if err != nil {
		log.Fatalf("Initservice failed loading test env %v", err)
	}

	svc, err := service.NewService(konf)
	if err != nil {
		log.Fatalf("Error initializing %v", err)
	}
	defer svc.Close()
	svc.DB.TruncateAll()

	before(svc)
	code := m.Run()
	after()

	//teardown
	os.Exit(code)
}
