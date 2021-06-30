package storage

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ip75/meteostation/config"
	"github.com/jmoiron/sqlx"
)

var PG Storage

type Storage struct {
	pg *sqlx.DB
}

func (s Storage) Init() error {

	fmt.Fprintln(os.Stdout, "storage: connecting to PostgreSQL database")
	d, err := sqlx.Open("postgres", config.C.PostgreSQL.DSN)
	if err != nil {
		return fmt.Errorf("storage: PostgreSQL connection error: %s", err)
	}
	d.SetMaxOpenConns(config.C.PostgreSQL.MaxOpenConnections)
	d.SetMaxIdleConns(config.C.PostgreSQL.MaxIdleConnections)
	for {
		if err := d.Ping(); err != nil {
			fmt.Fprintln(os.Stderr, "storage: ping PostgreSQL database error, will retry in 2s: ", err)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	PG = Storage{d}

	return err
}

func (s Storage) StoreSensorPoint(sensor SensorData) error {

	tx, err := s.pg.Beginx()
	if err != nil {
		return fmt.Errorf("storage: Unable to open transaction: %s", err)
	}

	_, err = tx.NamedExec("INSERT INTO sensor_data (date, temperature, pressure) VALUES (:date, :temperature, :pressure)", sensor)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("storage: transaction rollback error: %s", rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("storage: transaction commit error: %s", err)
	}
	return nil

}

func (s Storage) StoreSensorData(data []SensorData) error {

	if data == nil {
		return errors.New("storage: No data to store")
	}

	tx, err := s.pg.Beginx()
	if err != nil {
		return fmt.Errorf("storage: Unable to open transaction: %s", err)
	}

	_, err = tx.NamedExec("INSERT INTO sensor_data (date, temperature, pressure) VALUES (:date, :temperature, :pressure)", data)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("storage: transaction rollback error: %s", rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("storage: transaction commit error: %s", err)
	}
	return nil
}
