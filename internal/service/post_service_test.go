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
