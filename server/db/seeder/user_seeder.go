package seeder

import (
	"basic_server/server/model"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct {
	DB *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{DB: db}
}

type userData struct {
	Email    string
	Password string
	FullName string
}

func (userSeeder *UserSeeder) Run() {
	password, _ := bcrypt.GenerateFromPassword([]byte("test"), 14)
	users := map[int]userData{1: {
		Email:    "test@test.com",
		Password: string(password),
		FullName: "Test Test",
	}, 2: {
		Email:    "test1@test1.com",
		Password: string(password),
		FullName: "Test1 Test1",
	}}
	for key, value := range users {
		user := model.User{}
		userSeeder.DB.First(&user, key)
		if user.ID == 0 {
			user.ID = uint(key)
			user.Email = value.Email
			user.Password = value.Password
			user.FullName = value.FullName
			userSeeder.DB.Create(&user)
		}
	}
}
