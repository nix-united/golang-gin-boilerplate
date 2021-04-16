package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(200);UNIQUE"`
	Password string `gorm:"type:varchar(200);"`
	FullName string `gorm:"type:varchar(200);"`
}
