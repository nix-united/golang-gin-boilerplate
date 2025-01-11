package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/handler"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
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

	engine.POST("/posts", postHandler.SavePost)
	engine.GET("/posts", postHandler.GetPosts)
	engine.GET("/post/:id", postHandler.GetPostByID)
	engine.PUT("/post/:id", postHandler.UpdatePost)
	engine.DELETE("/post/:id", postHandler.DeletePost)

	return engine, postService
}

func TestPostHandler_GetPostByID(t *testing.T) {
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
		GetByID(100, &model.Post{}).
		DoAndReturn(func(i int, p *model.Post) *service.RestError {
			(*p) = post
			return nil
		})

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
		"content": "Content"
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
		CreatePost("Title", "Content", uint(101)).
		Return(&post, nil)

	httpRequest := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(rawCreatePostRequest))

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
		"content": "Content"
	}`

	assert.JSONEq(t, expectedResponse, string(responseBody))
}
