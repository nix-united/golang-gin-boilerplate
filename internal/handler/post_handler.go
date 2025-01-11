package handler

import (
	"net/http"
	"strconv"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"
	"github.com/nix-united/golang-gin-boilerplate/internal/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type postService interface {
	CreatePost(title, content string, userID uint) (*model.Post, *service.RestError)
	GetAll(posts *[]model.Post) *service.RestError
	GetByID(id int, post *model.Post) *service.RestError
	Create(post *model.Post) *service.RestError
	Save(post *model.Post) *service.RestError
	Delete(post *model.Post) *service.RestError
}

type PostHandler struct {
	postService postService
}

func NewPostHandler(postService postService) PostHandler {
	return PostHandler{
		postService: postService,
	}
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
func (h PostHandler) GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var post model.Post
	if err := h.postService.GetByID(id, &post); err != nil {
		response.ErrorResponse(c, err.Status, "Server error")
		return
	}

	if post.ID == 0 {
		response.ErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	response.SuccessResponse(c, response.GetPostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	})
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
func (h PostHandler) SavePost(c *gin.Context) {
	var createPostRequest request.CreatePostRequest
	if err := c.ShouldBindJSON(&createPostRequest); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
		return
	}

	claims := jwt.ExtractClaims(c)
	id := claims["id"].(float64)

	newPost, restError := h.postService.CreatePost(createPostRequest.Title, createPostRequest.Content, uint(id))
	if restError != nil {
		response.ErrorResponse(c, restError.Status, "Post can't be created")
		return
	}

	response.SuccessResponse(c, response.CreatePostResponse{
		ID:      newPost.ID,
		Title:   newPost.Title,
		Content: newPost.Content,
	})
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
func (h PostHandler) UpdatePost(c *gin.Context) {
	var updatePostRequest request.UpdatePostRequest
	if err := c.ShouldBindJSON(&updatePostRequest); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty")
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var post model.Post
	if err := h.postService.GetByID(id, &post); err != nil {
		response.ErrorResponse(c, err.Status, "Server error")
		return
	}

	if post.ID == 0 {
		response.ErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	post.Title = updatePostRequest.Title
	post.Content = updatePostRequest.Content

	if err := h.postService.Save(&post); err != nil {
		response.ErrorResponse(c, err.Status, "Data was not saved")
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
func (h PostHandler) GetPosts(c *gin.Context) {
	var posts []model.Post
	if err := h.postService.GetAll(&posts); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Server error")
		return
	}

	response.SuccessResponse(c, response.CreatePostsCollectionResponse(posts))
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
func (h PostHandler) DeletePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var post model.Post
	if err := h.postService.GetByID(id, &post); err != nil {
		response.ErrorResponse(c, err.Status, "Server error")
		return
	}

	if post.ID == 0 {
		response.ErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	if err := h.postService.Delete(&post); err != nil {
		response.ErrorResponse(c, err.Status, "Server error")
		return
	}

	response.SuccessResponse(c, "Post delete successfully")
}
