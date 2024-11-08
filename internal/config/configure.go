package config

type Config struct {
	DB   *DBConfig
	HTTP *HTTPConfig
}

func NewConfig() *Config {
	return &Config{
		DB:   LoadDBConfig(),
		HTTP: LoadHTTPConfig(),
	}
}
