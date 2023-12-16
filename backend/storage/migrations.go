package storage

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func (s Storage) PlayMigrations() {
	logrus.Info("Play migrations with DSN: ", s.ComposeDSN())
	m, err := migrate.New("file://storage/migrations", s.ComposeDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}
