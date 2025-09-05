package db

import (
	"database/sql"
	"fmt"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBConnection(cfg config.DBConfig) (gormDB *gorm.DB, sqlDB *sql.DB, err error) {
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	sqlDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, nil, fmt.Errorf("open db connection: %w", err)
	}

	gormDB, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{Logger: newLoggerAdapter()})
	if err != nil {
		return nil, nil, fmt.Errorf("open gorm session: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLife)

	return gormDB, sqlDB, nil
}
