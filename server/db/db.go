package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" //nolint
	"github.com/jinzhu/gorm"
)

func InitDB() *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	fmt.Println(dataSourceName)

	db, err := gorm.Open(os.Getenv("DB_DRIVER"), dataSourceName)

	if err != nil {
		panic(err.Error())
	}

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

	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLife) * time.Second)

	return db
}
