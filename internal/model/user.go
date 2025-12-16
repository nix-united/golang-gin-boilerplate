package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email    string
	Password string
	FullName string
}
