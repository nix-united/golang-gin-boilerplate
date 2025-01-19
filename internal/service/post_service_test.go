package service_test

import (
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestPostService_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := service.NewPostService(postRepository)

	expectedPostToCreate := &model.Post{
		Title:   "Title",
		Content: "Content",
		UserID:  100,
	}

	expectedCreatedPost := new(model.Post)
	*expectedCreatedPost = *expectedPostToCreate
	expectedCreatedPost.ID = 101

	postRepository.
		EXPECT().
		Create(expectedPostToCreate).
		DoAndReturn(func(p *model.Post) error {
			(*p) = *expectedCreatedPost
			return nil
		})

	post, err := postService.CreatePost("Title", "Content", 100)
	require.Nil(t, err)

	assert.Equal(t, expectedCreatedPost, post)
}

func TestPostService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := service.NewPostService(postRepository)

	postToCreate := &model.Post{
		Title:   "Title",
		Content: "Content",
		UserID:  100,
	}

	expectedCreatedPost := new(model.Post)
	*expectedCreatedPost = *postToCreate
	expectedCreatedPost.ID = 101

	postRepository.
		EXPECT().
		Create(postToCreate).
		DoAndReturn(func(p *model.Post) error {
			(*p) = *expectedCreatedPost
			return nil
		})

	err := postService.Create(postToCreate)
	assert.Nil(t, err)
}

func TestPostService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := service.NewPostService(postRepository)

	storedPosts := []model.Post{{
		Title:   "Title",
		Content: "Content",
		UserID:  100,
	}}

	postRepository.
		EXPECT().
		GetAll(gomock.Any()).
		DoAndReturn(func(p *[]model.Post) error {
			*p = storedPosts
			return nil
		})

	var posts []model.Post
	err := postService.GetAll(&posts)
	require.Nil(t, err)

	assert.Equal(t, storedPosts, posts)
}

func TestPostService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := service.NewPostService(postRepository)

	storedPost := &model.Post{
		Title:   "Title",
		Content: "Content",
		UserID:  100,
	}

	postRepository.
		EXPECT().
		GetByID(101).
		DoAndReturn(func(i int) (*model.Post, error) {
			return storedPost, nil
		})

	post, err := postService.GetByID(101)
	require.Nil(t, err)

	assert.Equal(t, storedPost, post)
}

func TestPostService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := service.NewPostService(postRepository)

	post := &model.Post{
		Model: gorm.Model{
			ID: 101,
		},
		Title:   "Title",
		Content: "Content",
		UserID:  102,
	}

	postRepository.
		EXPECT().
		Save(post).
		Return(nil)

	err := postService.Save(post)
	assert.Nil(t, err)
}

func TestPostService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := service.NewPostService(postRepository)

	post := &model.Post{
		Model: gorm.Model{
			ID: 101,
		},
		Title:   "Title",
		Content: "Content",
		UserID:  102,
	}

	postRepository.
		EXPECT().
		Delete(post).
		Return(nil)

	err := postService.Delete(post)
	assert.Nil(t, err)
}
