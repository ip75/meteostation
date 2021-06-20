package config

// Config defines the configuration structure.
type Config struct {
	General struct {
		LogLevel    int  `mapstructure:"log_level"`
		LogToSyslog bool `mapstructure:"log_to_syslog"`
	} `mapstructure:"general"`

	PostgreSQL struct {
		DSN                string `mapstructure:"dsn"`
		Automigrate        bool
		MaxOpenConnections int `mapstructure:"max_open_connections"`
		MaxIdleConnections int `mapstructure:"max_idle_connections"`
	} `mapstructure:"postgresql"`

	// #define REDIS_QUEUE "meteostation:bmp280"
	Redis struct {
		URL      string `mapstructure:"url"` // deprecated
		Password string `mapstructure:"password"`
		Database int    `mapstructure:"database"`
		Queue    string `mapstructure:"queue"`
	} `mapstructure:"redis"`
}

// C holds the global configuration.
var C Config

// Get returns the configuration.
func Get() *Config {
	return &C
}
