package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ip75/meteostation/config"
	"github.com/sirupsen/logrus"
)

var RDB RClient

type RClient struct {
	*redis.Client
}

type SensorData struct {
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Clock       int64   `json:"time"`
}

type SensorDataDatabase struct {
	Temperature float64   `db:"temperature"`
	Pressure    float64   `db:"pressure"`
	Date        time.Time `db:"dt"`
}

func (r RClient) Init() {
	fmt.Fprintln(os.Stdout, "storage: connect to Redis", config.C.Redis.URL)

	client := redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.URL,
		Password: config.C.Redis.Password, // no password set
		DB:       config.C.Redis.Database, // use default DB
	})

	RDB = RClient{client}
}

// pull data from redis queue with data from sensor
func (r RClient) PullPoint() SensorData {
	data, err := r.RPop(context.Background(), config.C.Redis.Queue).Result()
	if err != nil {
		logrus.Error("RPop error:", err)
		panic(err)
	}

	var sensor SensorData
	if err := json.Unmarshal([]byte(data), &sensor); err != nil {
		logrus.Error("unmarshal sensor data failed:", err)
	}

	return sensor
}

// pull data from redis queue with data from sensor
func (r RClient) Pull() []SensorDataDatabase {
	var result []SensorDataDatabase

	for i := 0; i < config.C.General.PoolSize; i++ {
		data, err := r.BRPop(context.Background(), 0, config.C.Redis.Queue).Result()
		if err != nil {
			logrus.Error("BRPop error:", err)
		}

		if len(data) < 2 {
			logrus.Errorf("BRPop error: no data from redis")
			return nil
		}

		point := SensorData{}
		if err = json.Unmarshal([]byte(data[1]), &point); err != nil {
			logrus.Errorf("unable to unmarshal: %s", err)
		}

		// This is not a timestamp this is clock from start of device.
		// so we overwrite it with current simestamp
		result = append(result, SensorDataDatabase{
			Temperature: point.Temperature,
			Pressure:    point.Pressure,
			Date:        time.Now(),
		})
	}

	logrus.Infof("storage: pull %d records from redis\n", len(result))

	return result
}
