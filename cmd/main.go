package main

import (
	"fmt"
	"github.com/sisukasco/henki/api"
	"github.com/sisukasco/henki/pkg/version"
	"log"
	"os"

	"github.com/sisukasco/commons/conf"

	"github.com/knadh/koanf"
	"github.com/spf13/pflag"
)

func main() {
	confx, err := conf.LoadConf(os.Args[1:], "HENKI_", func(flags *pflag.FlagSet) {
		flags.String("source", "./pkg/db/migrations", "Source folder for migrations")
	})

	if err != nil {
		log.Printf("Error loading conf %v", err)
		return
	}

	command := confx.Flags.Arg(0)

	switch command {
	case "serve":
		startServer(confx.Konf)
	case "version":
		showVersion()
	case "users":
		showUsers(confx.Konf)
	case "migrate":
		doMigration(confx)
	}
}

func startServer(konf *koanf.Koanf) error {
	log.Printf("Auth server starting. version %v build %v", version.Version, version.BuildTime)
	server, err := api.NewServer(konf)
	if err != nil {
		log.Fatal(err)
		return err
	}
	server.Start()
	return nil
}

func showVersion() {
	fmt.Printf("Henki version %s Build %s\n", version.Version, version.BuildTime)
}
