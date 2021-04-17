package db

import (
	"basic_server/config"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //nolint
	"github.com/jinzhu/gorm"
)

func InitDB(cfg *config.DBConfig) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := gorm.Open(cfg.Driver, dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	db.DB().SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.DB().SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLife) * time.Second)

	return db
}
