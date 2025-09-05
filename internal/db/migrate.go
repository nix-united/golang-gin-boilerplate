package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Migrate(dbConnection *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("set migrations dialect: %w", err)
	}

	if err := goose.Up(dbConnection, "migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
