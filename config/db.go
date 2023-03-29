package config

import (
	"log"
	"os"
	"strconv"
)

type DBConfig struct {
	User           string
	Password       string
	Driver         string
	Name           string
	Host           string
	Port           string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBConnMaxLife  int
}

func LoadDBConfig() *DBConfig {
	maxOpenConns, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		log.Fatal(err)
	}

	maxIdleConns, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		log.Fatal(err)
	}

	connMaxLife, err := strconv.Atoi(os.Getenv("DB_CONN_MAX_LIFE"))
	if err != nil {
		log.Fatal(err)
	}

	return &DBConfig{
		User:           os.Getenv("DB_USER"),
		Password:       os.Getenv("DB_PASSWORD"),
		Driver:         os.Getenv("DB_DRIVER"),
		Name:           os.Getenv("DB_NAME"),
		Host:           os.Getenv("DB_HOST"),
		Port:           os.Getenv("DB_PORT"),
		DBMaxOpenConns: maxOpenConns,
		DBMaxIdleConns: maxIdleConns,
		DBConnMaxLife:  connMaxLife,
	}
}
