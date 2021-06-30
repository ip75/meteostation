package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ip75/meteostation/config"
)

var RDB RClient

type RClient struct {
	Db *redis.Client
}

type SensorData struct {
	Temperature int       `json:"temperature" db:"temperature"`
	Pressure    int       `json:"pressure"    db:"pressure"`
	Date        time.Time `json:"date"        db:"date"`
}

var ctx = context.Background()

func (r RClient) Init() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.URL,
		Password: config.C.Redis.Password, // no password set
		DB:       config.C.Redis.Database, // use default DB
	})

	RDB = RClient{client}
}

// pull data from redis queue with data from sensor
func (r RClient) PullPoint() SensorData {
	data, err := r.Db.RPop(ctx, config.C.Redis.Queue).Result()

	if err != nil {
		panic(err)
	}

	var sensor SensorData
	json.Unmarshal([]byte(data), &sensor)

	return sensor
}

// pull data from redis queue with data from sensor
func (r RClient) Pull() []SensorData {

	data, err := r.Db.BRPop(ctx, 0, config.C.Redis.Queue).Result()

	if err != nil {
		panic(err)
	}

	var result []SensorData

	for _, s := range data {
		point := SensorData{}
		json.Unmarshal([]byte(s), &point)

		// This is not a timestamp this is clock from start of device.
		// so we overwrite it with current simestamp
		point.Date = time.Now()

		result = append(result, point)
	}

	return result
}
