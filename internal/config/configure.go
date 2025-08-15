package config

import "time"

type ApplicationConfig struct {
	ApplicationShutdownTimeout time.Duration `env:"APPLICATION_SHUTDOWN_TIMEOUT" envDefault:"5m"`

	DB         DBConfig
	HTTPServer HTTPServerConfig
}

type DBConfig struct {
	User           string        `env:"DB_USER"`
	Password       string        `env:"DB_PASSWORD"`
	Driver         string        `env:"DB_DRIVER"`
	Name           string        `env:"DB_NAME"`
	Host           string        `env:"DB_HOST"`
	Port           string        `env:"DB_PORT"`
	DBMaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" envDefault:"4"`
	DBMaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" envDefault:"4"`
	DBConnMaxLife  time.Duration `env:"DB_CONN_MAX_LIFE" envDefault:"5m"`
}

type HTTPServerConfig struct {
	Host              string        `env:"HOST"`
	Port              string        `env:"PORT" envDefault:"5500"`
	ReadHeaderTimeout time.Duration `env:"HTTP_SERVER_READ_HEADER_TIMEOUT" envDefault:"5m"`
	ReadTimeout       time.Duration `env:"HTTP_SERVER_READ_TIMEOUT" envDefault:"5m"`
	WriteTimeout      time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT" envDefault:"5m"`
}
