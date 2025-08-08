package post_test

import (
	"context"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/post"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestPostService_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

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
		Create(gomock.Any(), expectedPostToCreate).
		DoAndReturn(func(_ context.Context, p *model.Post) error {
			(*p) = *expectedCreatedPost
			return nil
		})

	post, err := postService.CreatePost(t.Context(), "Title", "Content", 100)
	require.Nil(t, err)

	assert.Equal(t, expectedCreatedPost, post)
}

func TestPostService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

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
		Create(gomock.Any(), postToCreate).
		DoAndReturn(func(_ context.Context, p *model.Post) error {
			(*p) = *expectedCreatedPost
			return nil
		})

	err := postService.Create(t.Context(), postToCreate)
	assert.Nil(t, err)
}

func TestPostService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	storedPosts := []model.Post{{
		Title:   "Title",
		Content: "Content",
		UserID:  100,
	}}

	postRepository.
		EXPECT().
		GetAll(gomock.Any()).
		Return(storedPosts, nil)

	posts, err := postService.GetAll(t.Context())
	require.Nil(t, err)

	assert.Equal(t, storedPosts, posts)
}

func TestPostService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	storedPost := &model.Post{
		Title:   "Title",
		Content: "Content",
		UserID:  100,
	}

	postRepository.
		EXPECT().
		GetByID(gomock.Any(), 101).
		Return(storedPost, nil)

	post, err := postService.GetByID(t.Context(), 101)
	require.Nil(t, err)

	assert.Equal(t, storedPost, post)
}

func TestPostService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

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
		Save(gomock.Any(), post).
		Return(nil)

	err := postService.Save(t.Context(), post)
	assert.Nil(t, err)
}

func TestPostService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

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
		Delete(gomock.Any(), post).
		Return(nil)

	err := postService.Delete(t.Context(), post)
	assert.Nil(t, err)
}
