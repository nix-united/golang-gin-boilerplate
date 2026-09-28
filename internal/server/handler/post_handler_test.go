package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"
	"github.com/nix-united/golang-gin-boilerplate/internal/server/handler"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func newPostHandler(t *testing.T, userID uint) (*gin.Engine, *MockpostService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	postService := NewMockpostService(ctrl)
	postHandler := handler.NewPostHandler(postService)

	engine := gin.New()
	gin.SetMode(gin.TestMode)

	engine.Use(func(c *gin.Context) {
		c.Set("JWT_PAYLOAD", jwt.MapClaims{"id": float64(userID)})
	})

	engine.POST("/posts", postHandler.CreatePost)
	engine.GET("/posts", postHandler.GetPosts)
	engine.GET("/posts/:id", postHandler.GetPostByID)
	engine.PUT("/posts/:id", postHandler.UpdatePost)
	engine.DELETE("/posts/:id", postHandler.DeletePost)

	return engine, postService
}

func TestPostHandler_CreatePost(t *testing.T) {
	const postID, userID = 100, 200
	now, err := time.Parse(time.DateOnly, "2020-01-02")
	require.NoError(t, err)

	createdPost := model.Post{
		Model: gorm.Model{
			ID:        postID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:  userID,
		Title:   "Title",
		Content: "Content",
	}

	createPostRequest := request.CreatePostRequest{
		BasicPost: request.BasicPost{
			Title:   createdPost.Title,
			Content: createdPost.Content,
		},
	}

	rawCreatePostRequest, err := json.Marshal(createPostRequest)
	require.NoError(t, err)

	engine, postService := newPostHandler(t, userID)

	postService.
		EXPECT().
		Create(gomock.Any(), domain.CreatePostRequest{
			UserID:  userID,
			Title:   createPostRequest.Title,
			Content: createPostRequest.Content,
		}).
		Return(&createdPost, nil)

	httpRequest := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(rawCreatePostRequest))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	httpResponse := recorder.Result()
	defer httpResponse.Body.Close()

	assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)

	responseBody, err := io.ReadAll(httpResponse.Body)
	require.NoError(t, err)

	var actualResponse response.PostResponse
	err = json.Unmarshal(responseBody, &actualResponse)
	require.NoError(t, err)

	expectedResponse := response.PostResponse{
		ID:        createdPost.ID,
		UserID:    createdPost.UserID,
		Title:     createdPost.Title,
		Content:   createdPost.Content,
		CreatedAt: createdPost.CreatedAt.Format(time.RFC3339),
		UpdatedAt: createdPost.UpdatedAt.Format(time.RFC3339),
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestPostHandler_GetPosts(t *testing.T) {
	const firstPostID, secondPostID, userID, totalPosts = 100, 101, 200, 300
	now, err := time.Parse(time.DateOnly, "2020-01-02")
	require.NoError(t, err)

	storedPosts := []model.Post{
		{
			Model: gorm.Model{
				ID:        firstPostID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:  userID,
			Title:   "Post 1 Title",
			Content: "Post 2 Content",
		},
		{
			Model: gorm.Model{
				ID:        secondPostID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:  userID,
			Title:   "Post 2 Title",
			Content: "Post 2 Content",
		},
	}

	filters := domain.PostFilters{Limit: 10}

	engine, postService := newPostHandler(t, userID)

	postService.
		EXPECT().
		Count(gomock.Any(), filters).
		Return(totalPosts, nil)

	postService.
		EXPECT().
		List(gomock.Any(), filters).
		Return(storedPosts, nil)

	httpRequest := httptest.NewRequest(http.MethodGet, "/posts", http.NoBody)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	httpResponse := recorder.Result()
	defer httpResponse.Body.Close()

	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

	responseBody, err := io.ReadAll(httpResponse.Body)
	require.NoError(t, err)

	var actualResponse response.CollectionResponse[response.PostResponse]
	err = json.Unmarshal(responseBody, &actualResponse)
	require.NoError(t, err)

	expectedResponse := response.CollectionResponse[response.PostResponse]{
		Data: []response.PostResponse{
			{
				ID:        storedPosts[0].ID,
				UserID:    storedPosts[0].UserID,
				Title:     storedPosts[0].Title,
				Content:   storedPosts[0].Content,
				CreatedAt: storedPosts[0].CreatedAt.Format(time.RFC3339),
				UpdatedAt: storedPosts[0].UpdatedAt.Format(time.RFC3339),
			},
			{
				ID:        storedPosts[1].ID,
				UserID:    storedPosts[1].UserID,
				Title:     storedPosts[1].Title,
				Content:   storedPosts[1].Content,
				CreatedAt: storedPosts[1].CreatedAt.Format(time.RFC3339),
				UpdatedAt: storedPosts[1].UpdatedAt.Format(time.RFC3339),
			},
		},
		Meta: response.Meta{
			Count:  2,
			Total:  totalPosts,
			Offset: 0,
			Limit:  10,
		},
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestPostHandler_GetPostByID(t *testing.T) {
	const postID, userID = 100, 200
	now, err := time.Parse(time.DateOnly, "2020-01-02")
	require.NoError(t, err)

	storedPost := &model.Post{
		Model: gorm.Model{
			ID:        postID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:  userID,
		Title:   "Title",
		Content: "Content",
	}

	engine, postService := newPostHandler(t, userID)

	postService.
		EXPECT().
		GetByID(gomock.Any(), storedPost.ID).
		Return(storedPost, nil)

	httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/posts/%d", postID), http.NoBody)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	httpResponse := recorder.Result()
	defer httpResponse.Body.Close()

	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

	responseBody, err := io.ReadAll(httpResponse.Body)
	require.NoError(t, err)

	var actualResponse response.PostResponse
	err = json.Unmarshal(responseBody, &actualResponse)
	require.NoError(t, err)

	expectedResponse := response.PostResponse{
		ID:        storedPost.ID,
		UserID:    storedPost.UserID,
		Title:     storedPost.Title,
		Content:   storedPost.Content,
		CreatedAt: storedPost.CreatedAt.Format(time.RFC3339),
		UpdatedAt: storedPost.UpdatedAt.Format(time.RFC3339),
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestPostHandler_UpdatePost(t *testing.T) {
	const userID = 101
	now, err := time.Parse(time.DateOnly, "2020-01-02")
	require.NoError(t, err)

	newPost := &model.Post{
		Model: gorm.Model{
			ID:        100,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:  userID,
		Title:   "New Title",
		Content: "New Content",
	}

	updatePostRequest := request.UpdatePostRequest{
		BasicPost: request.BasicPost{
			Title:   "New Title",
			Content: "New Content",
		},
	}

	rawUpdatePostRequest, err := json.Marshal(updatePostRequest)
	require.NoError(t, err)

	engine, postService := newPostHandler(t, userID)

	postService.
		EXPECT().
		UpdateByUser(gomock.Any(), domain.UpdatePostRequest{
			PostID:  newPost.ID,
			UserID:  newPost.UserID,
			Title:   newPost.Title,
			Content: newPost.Content,
		}).
		Return(newPost, nil)

	httpRequest := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/posts/%d", newPost.ID),
		bytes.NewReader(rawUpdatePostRequest),
	)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	httpResponse := recorder.Result()
	defer httpResponse.Body.Close()

	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

	responseBody, err := io.ReadAll(httpResponse.Body)
	require.NoError(t, err)

	var actualResponse response.PostResponse
	err = json.Unmarshal(responseBody, &actualResponse)
	require.NoError(t, err)

	expectedResponse := response.PostResponse{
		ID:        newPost.ID,
		UserID:    newPost.UserID,
		Title:     newPost.Title,
		Content:   newPost.Content,
		CreatedAt: newPost.CreatedAt.Format(time.RFC3339),
		UpdatedAt: newPost.UpdatedAt.Format(time.RFC3339),
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestPostHandler_DeletePost(t *testing.T) {
	const userID = uint(101)

	engine, postService := newPostHandler(t, userID)

	postService.
		EXPECT().
		DeleteByUser(gomock.Any(), userID, uint(100)).
		Return(nil)

	httpRequest := httptest.NewRequest(http.MethodDelete, "/posts/100", http.NoBody)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	httpResponse := recorder.Result()
	defer httpResponse.Body.Close()

	assert.Equal(t, http.StatusNoContent, httpResponse.StatusCode)

	responseBody, err := io.ReadAll(httpResponse.Body)
	require.NoError(t, err)

	assert.Empty(t, responseBody)
}
