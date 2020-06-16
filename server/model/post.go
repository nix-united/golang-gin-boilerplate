package model

import "github.com/jinzhu/gorm"

type Post struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint
	User    User `gorm:"foreignkey:UserID"`
}
