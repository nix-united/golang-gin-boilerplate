package config

type Config struct {
	DB   DB
	HTTP HTTP
}

type DB struct {
	User           string `env:"DB_USER"`
	Password       string `env:"DB_PASSWORD"`
	Driver         string `env:"DB_DRIVER"`
	Name           string `env:"DB_NAME"`
	Host           string `env:"DB_HOST"`
	Port           string `env:"DB_PORT"`
	DBMaxOpenConns int    `env:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns int    `env:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLife  int    `env:"DB_CONN_MAX_LIFE"`
}

type HTTP struct {
	Host       string `env:"HOST"`
	Port       string `env:"PORT"`
	ExposePort string `env:"EXPOSE_PORT"`
}
