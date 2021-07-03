package storage

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (s Storage) PlayMigrations() {

	fmt.Println("Play migrations with DSN: ", s.ComposeDSN())
	m, err := migrate.New("file://storage/migrations", s.ComposeDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
