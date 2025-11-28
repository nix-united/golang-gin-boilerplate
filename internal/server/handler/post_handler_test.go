package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func newPostHandler(t *testing.T) (*gin.Engine, *MockpostService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	postService := NewMockpostService(ctrl)
	postHandler := handler.NewPostHandler(postService)

	engine := gin.New()
	gin.SetMode(gin.TestMode)

	engine.Use(func(c *gin.Context) {
		c.Set("JWT_PAYLOAD", jwt.MapClaims{"id": float64(101)})
	})

	engine.POST("/posts", postHandler.CreatePost)
	engine.GET("/posts", postHandler.GetPosts)
	engine.GET("/post/:id", postHandler.GetPostByID)
	engine.PUT("/post/:id", postHandler.UpdatePost)
	engine.DELETE("/post/:id", postHandler.DeletePost)

	return engine, postService
}

func TestPostHandler_GetPostByID(t *testing.T) {
	engine, postService := newPostHandler(t)

	post := &model.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "Title",
		Content: "Content",
	}

	postService.
		EXPECT().
		GetByID(gomock.Any(), post.ID).
		Return(post, nil)

	httpRequest := httptest.NewRequest(http.MethodGet, "/post/100", http.NoBody)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	response := recorder.Result()
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	expectedResponse := `{
		"id": 100,
		"title": "Title",
		"content": "Content",
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}`

	assert.JSONEq(t, expectedResponse, string(responseBody))
}

func TestPostHandler_SavePost(t *testing.T) {
	engine, postService := newPostHandler(t)

	post := model.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "Title",
		Content: "Content",
	}

	createPostRequest := request.CreatePostRequest{
		BasicPost: &request.BasicPost{
			Title:   "Title",
			Content: "Content",
		},
	}

	rawCreatePostRequest, err := json.Marshal(createPostRequest)
	require.NoError(t, err)

	postService.
		EXPECT().
		Create(gomock.Any(), uint(101), "Title", "Content").
		Return(&post, nil)

	httpRequest := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(rawCreatePostRequest))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	response := recorder.Result()
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	expectedResponse := `{
		"id": 100,
		"title": "Title",
		"content": "Content",
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}`

	assert.JSONEq(t, expectedResponse, string(responseBody))
}

func TestPostHandler_UpdatePost(t *testing.T) {
	engine, postService := newPostHandler(t)

	post := &model.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "Title",
		Content: "Content",
		UserID:  101,
	}

	newPost := &model.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "New Title",
		Content: "New Content",
	}

	updatePostRequest := request.UpdatePostRequest{
		BasicPost: &request.BasicPost{
			Title:   "New Title",
			Content: "New Content",
		},
	}

	rawUpdatePostRequest, err := json.Marshal(updatePostRequest)
	require.NoError(t, err)

	postService.
		EXPECT().
		UpdateByUser(gomock.Any(), post.UserID, post.ID, "New Title", "New Content").
		Return(newPost, nil)

	httpRequest := httptest.NewRequest(http.MethodPut, "/post/100", bytes.NewReader(rawUpdatePostRequest))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	response := recorder.Result()
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	expectedResponse := `{
		"id": 100,
		"title": "New Title",
		"content": "New Content",
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}`

	assert.JSONEq(t, expectedResponse, string(responseBody))
}

func TestPostHandler_GetPosts(t *testing.T) {
	engine, postService := newPostHandler(t)

	post := model.Post{
		Model: gorm.Model{
			ID: 100,
		},
		Title:   "Title",
		Content: "Content",
	}

	postService.
		EXPECT().
		Count(gomock.Any()).
		Return(100, nil)

	postService.
		EXPECT().
		List(gomock.Any(), domain.PostFilters{Limit: 10}).
		Return([]model.Post{post}, nil)

	httpRequest := httptest.NewRequest(http.MethodGet, "/posts", http.NoBody)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	response := recorder.Result()
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	expectedResponse := `{
		"data": [
			{
				"id": 100,
				"title": "Title",
				"content": "Content",
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z"
			}
		],
		"meta": {
			"count": 1,
			"total": 100,
			"offset": 0,
			"limit": 10
		}
	}`

	assert.JSONEq(t, expectedResponse, string(responseBody))
}

func TestPostHandler_DeletePost(t *testing.T) {
	engine, postService := newPostHandler(t)

	postService.
		EXPECT().
		DeleteByUser(gomock.Any(), uint(101), uint(100)).
		Return(nil)

	httpRequest := httptest.NewRequest(http.MethodDelete, "/post/100", http.NoBody)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httpRequest)

	httpResponse := recorder.Result()
	defer httpResponse.Body.Close()

	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

	responseBody, err := io.ReadAll(httpResponse.Body)
	require.NoError(t, err)

	var gotMessageResponse response.MessageResponse
	err = json.Unmarshal(responseBody, &gotMessageResponse)
	require.NoError(t, err)

	wantMessageRespone := response.MessageResponse{
		Message: "Post was deleted successfully",
	}

	assert.Equal(t, wantMessageRespone, gotMessageResponse)
}
