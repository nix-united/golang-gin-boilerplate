package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

//go:generate go tool mockgen -source=$GOFILE -destination=post_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postService interface {
	Create(ctx context.Context, userID uint, title, content string) (*model.Post, error)
	GetByID(ctx context.Context, id uint) (*model.Post, error)
	List(ctx context.Context) ([]model.Post, error)
	UpdateByUser(ctx context.Context, userID, postID uint, title, content string) (*model.Post, error)
	DeleteByUser(ctx context.Context, userID, postID uint) error
}

type PostHandler struct {
	postService postService
}

func NewPostHandler(postService postService) *PostHandler {
	return &PostHandler{postService: postService}
}

// SavePost godoc
// @Summary Create post
// @Description Create post
// @ID posts-create
// @Tags Posts Actions
// @Accept json
// @Produce json
// @Param params body request.CreatePostRequest true "Post title and content"
// @Success 200 {string} response.CreatePostResponse
// @Failure 400 {string} string "Bad request"
// @Security ApiKeyAuth
// @Router /posts [post]
func (h *PostHandler) SavePost(c *gin.Context) {
	userID, ok := jwt.ExtractClaims(c)["id"].(float64)
	if !ok {
		response.ErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	var createPostRequest request.CreatePostRequest
	if err := c.ShouldBindJSON(&createPostRequest); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
		return
	}

	post, err := h.postService.Create(
		c.Request.Context(),
		uint(userID),
		createPostRequest.Title,
		createPostRequest.Content,
	)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Post can't be created")
		return
	}

	response.SuccessResponse(c, response.CreatePostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	})
}

// GetPostByID godoc
// @Summary Get post by id
// @Description Get post by id
// @ID get-post
// @Tags Posts Actions
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} response.GetPostResponse
// @Failure 401 {object} response.Error
// @Security ApiKeyAuth
// @Router /post/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	post, err := h.postService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, "Post not found")
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError, "Server error")
		return
	}

	response.SuccessResponse(c, response.GetPostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	})
}

// GetPosts godoc
// @Summary Get all posts
// @Description Get all posts of all users
// @ID get-posts
// @Tags Posts Actions
// @Produce json
// @Success 200 {object} response.CollectionResponse
// @Failure 401 {object} response.Error
// @Security ApiKeyAuth
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.postService.List(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Server error")
		return
	}

	response.SuccessResponse(c, response.CreatePostsCollectionResponse(posts))
}

// UpdatePost godoc
// @Summary Update post
// @Description Update post
// @ID posts-update
// @Tags Posts Actions
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param params body request.UpdatePostRequest true "Post title and content"
// @Success 200 {string} response.GetPostResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {object} response.Error
// @Security ApiKeyAuth
// @Router /post/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, ok := jwt.ExtractClaims(c)["id"].(float64)
	if !ok {
		response.ErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	var updatePostRequest request.UpdatePostRequest
	if err := c.ShouldBindJSON(&updatePostRequest); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	post, err := h.postService.UpdateByUser(
		c.Request.Context(),
		uint(userID),
		uint(postID),
		updatePostRequest.Title,
		updatePostRequest.Content,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			response.ErrorResponse(c, http.StatusNotFound, "Post not found")
		case errors.Is(err, domain.ErrForbidden):
			response.ErrorResponse(c, http.StatusForbidden, "Forbidden")
		default:
			response.ErrorResponse(c, http.StatusInternalServerError, "Server error")
		}
		return
	}

	response.SuccessResponse(c, response.GetPostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	})
}

// DeletePost godoc
// @Summary Delete post
// @Description Delete post
// @ID posts-delete
// @Tags Posts Actions
// @Param id path int true "Post ID"
// @Success 200 {string} string "Post deleted successfully"
// @Failure 404 {object} response.Error
// @Security ApiKeyAuth
// @Router /post/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, ok := jwt.ExtractClaims(c)["id"].(float64)
	if !ok {
		response.ErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Bad request")
		return
	}

	if err := h.postService.DeleteByUser(c.Request.Context(), uint(userID), uint(postID)); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			response.ErrorResponse(c, http.StatusNotFound, "Post not found")
		case errors.Is(err, domain.ErrForbidden):
			response.ErrorResponse(c, http.StatusForbidden, "Forbidden")
		default:
			response.ErrorResponse(c, http.StatusInternalServerError, "Server error")
		}
		return
	}

	response.SuccessResponse(c, "Post delete successfully")
}
