package handler

import (
	"basic_server/server/model"
	"basic_server/server/repository"
	"basic_server/server/request"
	"basic_server/server/response"
	"basic_server/server/service"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

// GetPost godoc
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
	postsRepository := repository.PostRepository{DB: handler.DB}
	post := model.Post{}
	id, _ := strconv.Atoi(context.Param("id"))
	postsRepository.GetByID(id, &post)

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

// CreatePost godoc
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

	if err := context.ShouldBind(&createPostRequest); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		return
	}

	claims := jwt.ExtractClaims(context)
	id := claims["id"].(float64)

	postsService := service.PostService{DB: handler.DB}
	newPost := postsService.CreatePost(createPostRequest.Title, createPostRequest.Content, uint(id))
	postsRepository := repository.PostRepository{DB: handler.DB}
	postsRepository.Create(&newPost)
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

	if err := context.ShouldBind(&updatePostRequest); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		return
	}

	postsRepository := repository.PostRepository{DB: handler.DB}
	post := model.Post{}
	id, _ := strconv.Atoi(context.Param("id"))
	postsRepository.GetByID(id, &post)

	if post.ID == 0 {
		response.ErrorResponse(context, http.StatusNotFound, "Post not found")
	}

	post.Title = updatePostRequest.Title
	post.Content = updatePostRequest.Content
	postsRepository.Save(&post)

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
	postsRepository := repository.PostRepository{DB: handler.DB}
	var posts []model.Post
	postsRepository.GetAll(&posts)
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
	postsRepository := repository.PostRepository{DB: handler.DB}
	post := model.Post{}
	id, _ := strconv.Atoi(context.Param("id"))
	postsRepository.GetByID(id, &post)

	if post.ID == 0 {
		response.ErrorResponse(context, http.StatusNotFound, "Post not found")
		return
	}

	postsRepository.Delete(&post)

	response.SuccessResponse(context, "Post delete successfully")
}
