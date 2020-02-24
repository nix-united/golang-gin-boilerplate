package seeder

import (
	"basic_server/server/model"
	"github.com/jinzhu/gorm"
)

type PostSeeder struct {
	DB *gorm.DB
}

func NewPostSeeder(db *gorm.DB) *PostSeeder {
	return &PostSeeder{DB: db}
}

type postData struct {
	Title   string
	Content string
	UserId  uint
}

func (postSeeder *PostSeeder) Run() {
	posts := map[int]postData{1: {
		Title:   "Post 1",
		Content: "Post1 Content",
		UserId:  1,
	}, 2: {
		Title:   "Post 2",
		Content: "Post2 Content",
		UserId:  2,
	}}
	for key, value := range posts {
		post := model.Post{}
		postSeeder.DB.First(&post, key)
		if post.ID == 0 {
			post.ID = uint(key)
			post.Title = value.Title
			post.Content = value.Content
			post.UserId = value.UserId
			postSeeder.DB.Create(&post)
		}
	}
}
