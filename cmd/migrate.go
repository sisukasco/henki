package main

import (
	"github.com/sisukasco/commons/conf"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func doMigration(confx *conf.Confx) {
	log.Printf("migration command ...")
	op := confx.Flags.Arg(1)
	if len(op) <= 0 {
		log.Fatal("migration operation is required. (up or reset) ")
		return
	}
	log.Printf("Migration command %s", op)
	source, _ := confx.Flags.GetString("source")
	if len(source) <= 0 {
		log.Fatal("The migration files source folder is required --source folder")
		return
	}
	log.Printf("Migration command. source  %s", source)
	dburl := confx.Konf.String("db.url")
	if len(dburl) <= 0 {
		log.Fatal("DB URL config is not set. It is required for the migration")
		return
	}
	migrateCmd(source, dburl, op)
}

func migrateCmd(sourceFolder string, dbURL string, op string) {
	sourceFolder = "file://" + sourceFolder
	m, err := migrate.New(
		sourceFolder,
		dbURL)
	if err != nil {
		log.Printf("Error creating migration %v", err)
		return
	}
	switch op {
	case "up":
		log.Printf("Migrating up ...")
		err = m.Up()
	case "reset":
		log.Printf("Migrating down ...")
		err = m.Drop()
	default:
		log.Printf("Migration command %s not identified", op)
	}

	if err != nil {
		log.Printf("Error Migrating %v", err)
	}
	os.Exit(0)
}
