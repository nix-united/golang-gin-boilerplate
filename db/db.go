package db

import (
	"basic_server/config"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //nolint
	"github.com/jinzhu/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := gorm.Open(cfg.DB.Driver, dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	db.DB().SetMaxOpenConns(cfg.DB.DBMaxOpenConns)
	db.DB().SetMaxIdleConns(cfg.DB.DBMaxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(cfg.DB.DBConnMaxLife) * time.Second)

	return db
}
