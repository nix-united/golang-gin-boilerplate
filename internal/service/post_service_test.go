package service_test

import (
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
