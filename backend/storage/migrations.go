package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func PlayMigrations() {

	pgHost, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		pgHost = "pg"
	}
	pgDB, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		pgDB = "meteostation"
	}
	pgUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		pgUser = "meteostation"
	}
	pgPassword, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		pgPassword = "meteostation"
	}

	m, err := migrate.New(
		"file://backend/storage/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgDB))
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
