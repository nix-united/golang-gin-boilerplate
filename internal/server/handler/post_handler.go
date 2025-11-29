package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nix-united/golang-gin-boilerplate/internal/domain"
	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"

	safecast "github.com/ccoveille/go-safecast"
	"github.com/gin-gonic/gin"
)

const defaultPostLimit = 10

//go:generate go tool mockgen -source=$GOFILE -destination=post_handler_mock_test.go -package=${GOPACKAGE}_test -typed=true

type postService interface {
	Create(ctx context.Context, createPostRequest domain.CreatePostRequest) (*model.Post, error)
	Count(ctx context.Context, filers domain.PostFilters) (int64, error)
	List(ctx context.Context, filers domain.PostFilters) ([]model.Post, error)
	GetByID(ctx context.Context, postID uint) (*model.Post, error)
	UpdateByUser(ctx context.Context, updatePostRequest domain.UpdatePostRequest) (*model.Post, error)
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
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Error(fmt.Errorf("get user id from context: %w", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeInternalServerError,
			"Oops, something went wrong...",
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

	if err := createPostRequest.Validate(); err != nil {
		c.Error(fmt.Errorf("validate: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	post, err := h.postService.Create(c.Request.Context(), domain.CreatePostRequest{
		UserID:  userID,
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
	})
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

// GetPosts godoc
// @Summary Get all posts
// @ID getPosts
// @Tags Posts Actions
// @Produce json
// @Param limit query int false "Limit" minimum(1) maximum(100)
// @Param offset query int false "Offset" minimum(0)
// @Param user_id query string false "User ID"
// @Param title query string false "Title"
// @Success 200 {object} response.PostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	filters, err := h.parseFilters(c)
	if err != nil {
		c.Error(fmt.Errorf("parse post filters: %w", err))
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(
			response.CodeBadRequest,
			"Invalid request",
		))
		return
	}

	total, err := h.postService.Count(c.Request.Context(), filters)
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
// @Router /posts/{id} [get]
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
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Error(fmt.Errorf("get user id from context: %w", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeInternalServerError,
			"Oops, something went wrong...",
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

	post, err := h.postService.UpdateByUser(c.Request.Context(), domain.UpdatePostRequest{
		PostID:  postID,
		UserID:  userID,
		Title:   updatePostRequest.Title,
		Content: updatePostRequest.Content,
	})
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
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Error(fmt.Errorf("get user id from context: %w", err))
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
			response.CodeInternalServerError,
			"Oops, something went wrong...",
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

func (h *PostHandler) parseFilters(c *gin.Context) (domain.PostFilters, error) {
	filters := domain.PostFilters{Limit: defaultPostLimit, Title: c.Query("title")}

	if limitParam := c.Query("limit"); limitParam != "" {
		limit, err := strconv.ParseInt(limitParam, 10, 64)
		if err != nil {
			return domain.PostFilters{}, fmt.Errorf("prase limit query param: %w", err)
		}

		filters.Limit = limit
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		offset, err := strconv.ParseInt(offsetParam, 10, 64)
		if err != nil {
			return domain.PostFilters{}, fmt.Errorf("prase offset query param: %w", err)
		}

		filters.Offset = offset
	}

	if userIDQuery := c.Query("user_id"); userIDQuery != "" {
		userID, err := strconv.ParseUint(userIDQuery, 10, 64)
		if err != nil {
			return domain.PostFilters{}, fmt.Errorf("prase user_id query param: %w", err)
		}

		parsedUserID, err := safecast.ToUint(userID)
		if err != nil {
			return domain.PostFilters{}, fmt.Errorf("convert user id to uint: %w", err)
		}

		filters.UserID = parsedUserID
	}

	if err := filters.Validate(); err != nil {
		return domain.PostFilters{}, fmt.Errorf("validate filters: %w", err)
	}

	return filters, nil
}
