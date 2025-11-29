package post_test

import (
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/post"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestService_Create(t *testing.T) {
	createPostRequest := domain.CreatePostRequest{
		UserID:  100,
		Title:   "Title",
		Content: "Content",
	}

	postToCreate := &model.Post{
		UserID:  createPostRequest.UserID,
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
	}

	createdPost := &model.Post{
		Model:   gorm.Model{ID: 100},
		UserID:  createPostRequest.UserID,
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
	}

	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	postRepository.
		EXPECT().
		Create(gomock.Any(), postToCreate).
		Return(createdPost, nil)

	post, err := postService.Create(t.Context(), createPostRequest)
	require.NoError(t, err)

	assert.Equal(t, createdPost, post)
}

func TestService_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	const wantCount = int64(100)

	postRepository.
		EXPECT().
		Count(gomock.Any()).
		Return(wantCount, nil)

	gotCount, err := postService.Count(t.Context())
	require.NoError(t, err)

	assert.Equal(t, wantCount, gotCount)
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	storedPosts := []model.Post{{
		Title:   "Title",
		Content: "Content",
		UserID:  100,
	}}

	filters := domain.PostFilters{
		Offset: 1,
		Limit:  10,
	}

	postRepository.
		EXPECT().
		List(gomock.Any(), filters).
		Return(storedPosts, nil)

	posts, err := postService.List(t.Context(), filters)
	require.NoError(t, err)

	assert.Equal(t, storedPosts, posts)
}

func TestService_GetByID(t *testing.T) {
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

func TestService_UpdateByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	postRepository := NewMockpostRepository(ctrl)
	postService := post.NewService(postRepository)

	storedPost := &model.Post{
		Model: gorm.Model{
			ID: 101,
		},
		Title:   "Title",
		Content: "Content",
		UserID:  102,
	}

	updatedPost := &model.Post{
		Model: gorm.Model{
			ID: storedPost.ID,
		},
		Title:   "New Title",
		Content: "New Content",
		UserID:  storedPost.UserID,
	}

	postRepository.
		EXPECT().GetByID(gomock.Any(), storedPost.ID).
		Return(storedPost, nil)

	postRepository.
		EXPECT().
		Update(gomock.Any(), updatedPost).
		Return(nil)

	gotPost, err := postService.UpdateByUser(t.Context(), domain.UpdatePostRequest{
		PostID:  updatedPost.ID,
		UserID:  updatedPost.UserID,
		Title:   updatedPost.Title,
		Content: updatedPost.Content,
	})
	require.NoError(t, err)

	assert.Equal(t, updatedPost, gotPost)
}

func TestService_DeleteByUser(t *testing.T) {
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
		GetByID(gomock.Any(), post.ID).
		Return(post, nil)

	postRepository.
		EXPECT().
		Delete(gomock.Any(), post).
		Return(nil)

	err := postService.DeleteByUser(t.Context(), post.UserID, post.ID)
	assert.NoError(t, err)
}
