package storage

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ip75/meteostation/config"
	meteostation "github.com/ip75/meteostation/proto/api"
	"github.com/jmoiron/sqlx"
)

var PG Storage

type Storage struct {
	pg *sqlx.DB
}

func (s Storage) ComposeDSN() string {

	pgHost, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return config.C.PostgreSQL.DSN
	}
	pgDB, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return config.C.PostgreSQL.DSN
	}
	pgUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return config.C.PostgreSQL.DSN
	}
	pgPassword, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return config.C.PostgreSQL.DSN
	}

	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgDB)
}

func (s Storage) Init() error {

	fmt.Fprintln(os.Stdout, "storage: connecting to PostgreSQL database...")
	d, err := sqlx.Open("postgres", s.ComposeDSN())
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
	fmt.Fprintln(os.Stdout, "storage: connected to PostgreSQL database")

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

func (s Storage) StoreSensorData(data []SensorDataDatabase) error {

	if data == nil {
		return errors.New("storage: No data to store")
	}

	tx, err := s.pg.Beginx()
	if err != nil {
		return fmt.Errorf("storage: Unable to open transaction: %s", err)
	}

	_, err = tx.NamedExec("INSERT INTO meteodata (dt, temperature, pressure) VALUES (:dt, :temperature, :pressure)", data)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("storage: transaction rollback error: %s", rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("storage: transaction commit error: %s", err)
	}

	fmt.Printf("storage: dump %d records to database\n", len(data))

	return nil
}

type MeteoDataEntity struct {
	dt          time.Time `db:"dt"`
	temperature int64     `db:"temperature"`
	pressure    int64     `db:"pressure"`
	altitude    int64     `db:"altitude"`
}

func (s Storage) GetMeteoData(filter *meteostation.Filter) (meteostation.MeteoData, error) {

	data := MeteoDataEntity{
		dt:          time.Now(),
		temperature: 0,
		pressure:    0,
		altitude:    0,
	}

	var from time.Time = time.Now().Add(time.Hour * 24 * -31)
	var to time.Time = time.Now()
	var granularity int64 = 1000

	if filter != nil {
		from = filter.From.AsTime()
		to = filter.To.AsTime()
		granularity = filter.Granularity
	}

	s.pg.Get(&data, "SELECT * FROM meteodata", from, to, granularity)

	return meteostation.MeteoData{TotalCount: 0}, nil
}
