# Meteostation
This is a system that collects data from ESP32 + BPM280 sensor (temperature and pressure) to database to show changing dynamic of these parameters.

# Install

- Programm for ESP32 microcontroller: `arduino\Meteostation\Meteostation.ino`
- ESP32 module connects to WiFi AP and pushes data to redis database located at `$REDIS_HOST:$REDIS_PORT`. The main thing to deploy meteostation is to make DNS server to resolve host `$REDIS_HOST` as a host where you run containers by `docker-compose`.
- All settings for image are located in config file `backend\.meteostation.json` 

### Web config for ESP32 is in development...



## To create migration if you are going to extend schema

`migrate -source file://./backend/storage/migrations -database postgres://localhost:5432/meteostaion create -dir backend/storage/migrations -ext sql initialize`
