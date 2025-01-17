package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.DBConfig) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	sqlDB, err := sql.Open(cfg.Driver, dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLife) * time.Second)

	return gormDB
}
