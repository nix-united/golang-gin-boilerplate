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
func (handler PostHandler) GetPostByID(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	post := model.Post{}
	if err := handler.postService.GetByID(id, &post); err != nil {
		response.ErrorResponse(context, err.Status, "Server error")
		return
	}

	if post.ID == 0 {
		response.ErrorResponse(context, http.StatusNotFound, "Post not found")
		return
	}

	response.SuccessResponse(context, response.GetPostResponse{
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
func (handler PostHandler) SavePost(context *gin.Context) {
	var createPostRequest request.CreatePostRequest
	if err := context.ShouldBindJSON(&createPostRequest); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		return
	}

	claims := jwt.ExtractClaims(context)
	id := claims["id"].(float64)

	newPost, restError := handler.postService.CreatePost(createPostRequest.Title, createPostRequest.Content, uint(id))
	if restError != nil {
		response.ErrorResponse(context, restError.Status, "Post can't be created")
		return
	}

	response.SuccessResponse(context, response.CreatePostResponse{
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
func (handler PostHandler) UpdatePost(context *gin.Context) {
	var updatePostRequest request.UpdatePostRequest

	if err := context.ShouldBindJSON(&updatePostRequest); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		return
	}

	id, _ := strconv.Atoi(context.Param("id"))

	post := model.Post{}
	if err := handler.postService.GetByID(id, &post); err != nil {
		response.ErrorResponse(context, err.Status, "Server error")
		return
	}

	if post.ID == 0 {
		response.ErrorResponse(context, http.StatusNotFound, "Post not found")
		return
	}

	post.Title = updatePostRequest.Title
	post.Content = updatePostRequest.Content
	if err := handler.postService.Save(&post); err != nil {
		response.ErrorResponse(context, err.Status, "Data was not saved")
		return
	}

	response.SuccessResponse(context, response.GetPostResponse{
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
func (handler PostHandler) GetPosts(context *gin.Context) {
	var posts []model.Post
	if err := handler.postService.GetAll(&posts); err != nil {
		response.ErrorResponse(context, http.StatusInternalServerError, "Server error")
		return
	}
	response.SuccessResponse(context, response.CreatePostsCollectionResponse(posts))
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
func (handler PostHandler) DeletePost(context *gin.Context) {
	post := model.Post{}
	id, _ := strconv.Atoi(context.Param("id"))
	if err := handler.postService.GetByID(id, &post); err != nil {
		response.ErrorResponse(context, err.Status, "Server error")
		return
	}

	if post.ID == 0 {
		response.ErrorResponse(context, http.StatusNotFound, "Post not found")
		return
	}

	if err := handler.postService.Delete(&post); err != nil {
		response.ErrorResponse(context, err.Status, "Server error")
		return
	}

	response.SuccessResponse(context, "Post delete successfully")
}
