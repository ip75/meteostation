package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ip75/meteostation/config"
)

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
	r.Db = redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.URL,
		Password: config.C.Redis.Password, // no password set
		DB:       config.C.Redis.Database, // use default DB
	})
}

// pull data from redis queue with data from sensor
func (r RClient) Pull() SensorData {
	data := r.Db.RPop(ctx, config.C.Redis.Queue)
	if err := data.Err(); err != nil {
		panic(err)
	}

	res := SensorData{}
	json.Unmarshal([]byte(data.Val()), &res)

	return res
}
