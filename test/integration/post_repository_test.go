package integration_test

import (
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepository(t *testing.T) {
	userRepository := repository.NewUserRepository(gormDB)
	postRepository := repository.NewPostRepository(gormDB)

	user := &model.User{
		Email:    "test_post_repository_user@example.com",
		Password: "test_post_repository_user_password",
		FullName: "test_post_repository_user_fullname",
	}

	t.Run("Create user for posts", func(t *testing.T) {
		err := userRepository.Create(t.Context(), user)
		require.NoError(t, err)
	})

	post := &model.Post{
		Title:   "test_post_repository_title",
		Content: "test_post_repository_content",
		UserID:  user.ID,
	}

	t.Run("It should create a post", func(t *testing.T) {
		err := postRepository.Create(t.Context(), post)
		require.NoError(t, err)
	})

	t.Run("It should get post by ID", func(t *testing.T) {
		gotPost, err := postRepository.GetByID(t.Context(), post.ID)
		require.NoError(t, err)

		assert.Equal(t, post.Title, gotPost.Title)
	})

	t.Run("It should return ErrNotFound error when post with such ID not found", func(t *testing.T) {
		_, err := postRepository.GetByID(t.Context(), 999999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("It should fetch all posts", func(t *testing.T) {
		gotPosts, err := postRepository.List(t.Context())
		require.NoError(t, err)

		require.Len(t, gotPosts, 1)
		assert.Equal(t, post.Title, gotPosts[0].Title)
	})

	t.Run("It should update existing post", func(t *testing.T) {
		post.Title = "test_post_repository_title_updated"

		err := postRepository.Update(t.Context(), post)
		require.NoError(t, err)

		gotPost, err := postRepository.GetByID(t.Context(), post.ID)
		require.NoError(t, err)

		assert.Equal(t, post.Title, gotPost.Title)
	})

	t.Run("It should soft delete existing post", func(t *testing.T) {
		err := postRepository.Delete(t.Context(), post)
		require.NoError(t, err)
	})
}
