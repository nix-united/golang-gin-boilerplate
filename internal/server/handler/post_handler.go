package handler

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"

	jwt "github.com/appleboy/gin-jwt/v3"
	safecast "github.com/ccoveille/go-safecast"
	"github.com/gin-gonic/gin"
)

const (
	defaultPostLimit = 10
	maxPostLimit     = 100
)

//go:generate go tool mockgen -source=$GOFILE -destination=post_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postService interface {
	Create(ctx context.Context, userID uint, title, content string) (*model.Post, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, filers domain.PostFilters) ([]model.Post, error)
	GetByID(ctx context.Context, id uint) (*model.Post, error)
	UpdateByUser(ctx context.Context, userID, postID uint, title, content string) (*model.Post, error)
	DeleteByUser(ctx context.Context, userID, postID uint) error
}

type PostHandler struct {
	postService postService
}

func NewPostHandler(postService postService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePost godoc
// @Summary Create post
// @ID createPost
// @Tags Posts Actions
// @Accept json
// @Produce json
// @Param params body request.CreatePostRequest true "Post title and content"
// @Success 201 {object} response.PostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	parsedUserID, ok := claims["id"].(float64)
	if !ok {
		c.Error(fmt.Errorf("missing user id in claims: %v", claims))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	userID, err := safecast.ToUint(parsedUserID)
	if err != nil {
		c.Error(fmt.Errorf("convert user id to uint: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	var createPostRequest request.CreatePostRequest
	if err := c.ShouldBindJSON(&createPostRequest); err != nil {
		c.Error(fmt.Errorf("bind: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	post, err := h.postService.Create(
		c.Request.Context(),
		userID,
		createPostRequest.Title,
		createPostRequest.Content,
	)
	if err != nil {
		c.Error(fmt.Errorf("create post: %w", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeInternalServerError,
			"Oops, something went wrong...",
		))
		return
	}

	c.JSON(http.StatusCreated, response.NewPostResponse(post))
}

// GetPostByID godoc
// @Summary Get post by ID
// @ID getPostById
// @Tags Posts Actions
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} response.PostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /post/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	parsedPostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("parse post id: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	postID, err := safecast.ToUint(parsedPostID)
	if err != nil {
		c.Error(fmt.Errorf("convert post id to uint: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	post, err := h.postService.GetByID(c.Request.Context(), postID)
	if err != nil {
		c.Error(fmt.Errorf("get post by id: %w", err))
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(
				response.CodeBadRequest,
				"Post not found",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeInternalServerError,
			"Oops, something went wrong...",
		))
		return
	}

	c.JSON(http.StatusOK, response.NewPostResponse(post))
}

// GetPosts godoc
// @Summary Get all posts
// @ID getPosts
// @Tags Posts Actions
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.PostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	var (
		limit int64
		err   error
	)

	if param := c.Param("limit"); param != "" {
		limit, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			c.Error(fmt.Errorf("parse limit: %w", err))
			c.JSON(http.StatusBadRequest, response.NewErrorResponse(
				response.CodeBadRequest,
				"Invalid request",
			))
			return
		}
	}

	if limit > maxPostLimit {
		c.Error(fmt.Errorf("extra large posts limit: %d", limit))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	var offset int64
	if param := c.Param("offset"); param != "" {
		offset, err = strconv.ParseInt(param, 10, 64)
		if err != nil {
			c.Error(fmt.Errorf("parse offset: %w", err))
			c.JSON(http.StatusBadRequest, response.NewErrorResponse(
				response.CodeBadRequest,
				"Invalid request",
			))
			return
		}
	}

	filters := domain.PostFilters{
		Offset: offset,
		Limit:  cmp.Or(limit, defaultPostLimit),
	}

	total, err := h.postService.Count(c.Request.Context())
	if err != nil {
		c.Error(fmt.Errorf("count posts: %w", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeInternalServerError,
			"Oops, something went wrong...",
		))
		return
	}
	if total == 0 {
		c.JSON(http.StatusOK, response.NewPostCollectionResponse(nil, 0, filters.Offset, filters.Limit))
		return
	}

	posts, err := h.postService.List(c.Request.Context(), filters)
	if err != nil {
		c.Error(fmt.Errorf("list posts: %w", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeInternalServerError,
			"Oops, something went wrong...",
		))
		return
	}
	c.JSON(http.StatusOK, response.NewPostCollectionResponse(posts, total, filters.Offset, filters.Limit))
}

// UpdatePost godoc
// @Summary Update post
// @ID updatePost
// @Tags Posts Actions
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param params body request.UpdatePostRequest true "Post title and content"
// @Success 200 {string} response.PostResponse
// @Failure 400 {string} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /post/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	parsedUserID, ok := claims["id"].(float64)
	if !ok {
		c.Error(fmt.Errorf("missing user id in claims: %v", claims))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	userID, err := safecast.ToUint(parsedUserID)
	if err != nil {
		c.Error(fmt.Errorf("convert user id to uint: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	parsedPostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("parse post id: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	postID, err := safecast.ToUint(parsedPostID)
	if err != nil {
		c.Error(fmt.Errorf("convert post id to uint: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	var updatePostRequest request.UpdatePostRequest
	if err := c.ShouldBindJSON(&updatePostRequest); err != nil {
		c.Error(fmt.Errorf("bind: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	post, err := h.postService.UpdateByUser(
		c.Request.Context(),
		userID,
		postID,
		updatePostRequest.Title,
		updatePostRequest.Content,
	)
	if err != nil {
		c.Error(fmt.Errorf("update post by user: %w", err))
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, response.NewErrorResponse(
				response.CodeNotFound,
				"Post not found",
			))
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, response.NewErrorResponse(
				response.CodeAccessDenied,
				"Access denied",
			))
		default:
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
				response.CodeInternalServerError,
				"Oops, something went wrong...",
			))
		}
		return
	}

	c.JSON(http.StatusOK, response.NewPostResponse(post))
}

// DeletePost godoc
// @Summary Delete post
// @ID detelePost
// @Tags Posts Actions
// @Param id path int true "Post ID"
// @Success 200 {string} response.MessageResponse
// @Failure 400 {string} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /post/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	parsedUserID, ok := claims["id"].(float64)
	if !ok {
		c.Error(fmt.Errorf("missing user id in claims: %v", claims))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	userID, err := safecast.ToUint(parsedUserID)
	if err != nil {
		c.Error(fmt.Errorf("convert user id to uint: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	parsedPostID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("parse post id: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	postID, err := safecast.ToUint(parsedPostID)
	if err != nil {
		c.Error(fmt.Errorf("convert post id to uint: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	if err := h.postService.DeleteByUser(c.Request.Context(), userID, postID); err != nil {
		c.Error(fmt.Errorf("delete post by user: %w", err))
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, response.NewErrorResponse(
				response.CodeNotFound,
				"Post not found",
			))
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, response.NewErrorResponse(
				response.CodeAccessDenied,
				"Access denied",
			))
		default:
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
				response.CodeInternalServerError,
				"Oops, something went wrong...",
			))
		}
		return
	}

	c.JSON(http.StatusOK, response.NewMessageResponse("Post was deleted successfully"))
}
