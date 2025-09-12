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

func TestPostService_Create(t *testing.T) {
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

	post, err := postService.Create(t.Context(), 100, "Title", "Content")
	require.NoError(t, err)

	assert.Equal(t, expectedCreatedPost, post)
}

func TestPostService_List(t *testing.T) {
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
		List(gomock.Any()).
		Return(storedPosts, nil)

	posts, err := postService.List(t.Context())
	require.NoError(t, err)

	assert.Equal(t, storedPosts, posts)
}

func TestPostService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	storedPost := &model.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "Title",
		Content: "Content",
		UserID:  101,
	}

	postRepository.
		EXPECT().
		GetByID(gomock.Any(), storedPost.ID).
		Return(storedPost, nil)

	post, err := postService.GetByID(t.Context(), 100)
	require.NoError(t, err)

	assert.Equal(t, storedPost, post)
}

func TestPostService_UpdateByUser(t *testing.T) {
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

	newPost := &model.Post{
		Model: gorm.Model{
			ID: 101,
		},
		Title:   "New Title",
		Content: "New Content",
		UserID:  102,
	}

	postRepository.
		EXPECT().GetByID(gomock.Any(), post.ID).
		Return(post, nil)

	postRepository.
		EXPECT().
		Update(gomock.Any(), newPost).
		Return(nil)

	gotPost, err := postService.UpdateByUser(t.Context(), post.UserID, post.ID, "New Title", "New Content")
	require.NoError(t, err)

	assert.Equal(t, newPost, gotPost)
}

// func TestPostService_Delete(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	postRepository := NewMockpostRepository(ctrl)
// 	postService := post.NewService(postRepository)

// 	post := &model.Post{
// 		Model: gorm.Model{
// 			ID: 101,
// 		},
// 		Title:   "Title",
// 		Content: "Content",
// 		UserID:  102,
// 	}

// 	postRepository.
// 		EXPECT().
// 		Delete(gomock.Any(), post).
// 		Return(nil)

// 	err := postService.Delete(t.Context(), post)
// 	assert.Nil(t, err)
// }
