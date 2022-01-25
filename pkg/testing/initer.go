package testing

import (
	"log"
	"os"
	"testing"

	"github.com/sisukasco/commons/conf"
	"github.com/sisukasco/henki/pkg/service"
	"github.com/spf13/pflag"
)

func InitService(m *testing.M, confPath string, before func(s *service.Service), after func()) {

	args := []string{"--conf", confPath}

	//setup
	var err error
	confx, err := conf.LoadConf(args, "HENKI_", func(flags *pflag.FlagSet) {})

	if err != nil {
		log.Fatalf("Initservice failed loading test conf %v", err)
	}

	svc, err := service.NewService(confx.Konf)
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
