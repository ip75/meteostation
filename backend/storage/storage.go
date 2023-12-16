package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ip75/meteostation/config"
	meteostation "github.com/ip75/meteostation/proto/api"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		return fmt.Errorf("storage: PostgreSQL connection error: %w", err)
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
		return fmt.Errorf("storage: Unable to open transaction: %w", err)
	}

	_, err = tx.NamedExec("INSERT INTO sensor_data (date, temperature, pressure) VALUES (:date, :temperature, :pressure)", sensor)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("storage: transaction rollback error: %w", rbErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("storage: transaction commit error: %w", err)
	}
	return nil
}

func (s Storage) StoreSensorData(data []SensorDataDatabase) error {
	if data == nil {
		return errors.New("storage: No data to store")
	}

	tx, err := s.pg.Beginx()
	if err != nil {
		return fmt.Errorf("storage: Unable to open transaction: %w", err)
	}

	_, err = tx.NamedExec("INSERT INTO meteodata (dt, temperature, pressure) VALUES (:dt, :temperature, :pressure)", data)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("storage: transaction rollback error: %w", rbErr)
		}

		return fmt.Errorf("storage: insert meteodata error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("storage: transaction commit error: %w", err)
	}

	logrus.Infof("%s : storage : dump %d records to database", time.Now(), len(data))

	return nil
}

type MeteoDataEntity struct {
	Dt          time.Time       `db:"dt"`
	Temperature float64         `db:"temperature"`
	Pressure    float64         `db:"pressure"`
	Altitude    sql.NullFloat64 `db:"altitude"`
}

func (s Storage) GetMeteoData(filter *meteostation.Filter) (*meteostation.MeteoData, error) {
	var data []MeteoDataEntity

	from := time.Now().Add(time.Hour * 24 * -31)
	to := time.Now()
	granularity := int64(1)

	if filter != nil {
		from = filter.From.AsTime()
		to = filter.To.AsTime()
		granularity = filter.Granularity
	}

	err := s.pg.Select(&data, `
			SELECT 
				dt,
				temperature,
				pressure,
				altitude 
			FROM meteodata m
			WHERE
				dt BETWEEN $1 AND $2
				AND id % $3 = 0`,
		from,
		to,
		granularity,
	)
	if err != nil {
		return nil, fmt.Errorf("get meteodata:%w", err)
	}

	sd := make([]*meteostation.SensorData, 0, len(data))
	for _, d := range data {
		sd = append(sd, &meteostation.SensorData{
			Temperature: d.Temperature,
			Pressure:    d.Pressure,
			Altitude:    d.Altitude.Float64,
			MeasureTime: timestamppb.New(d.Dt),
		})
	}

	return &meteostation.MeteoData{TotalCount: uint64(len(sd)), Data: sd}, nil
}
