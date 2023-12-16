package storage

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func (s Storage) PlayMigrations() {
	logrus.Info("Play migrations with DSN: ", s.ComposeDSN())
	m, err := migrate.New("file://storage/migrations", s.ComposeDSN())
	if err != nil {
		logrus.Fatal("create migration failed: ", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatal("play migration failed: ", err)
	}
}
