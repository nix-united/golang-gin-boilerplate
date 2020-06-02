package handler

import (
	"basic_server/server/model"
	"basic_server/server/repository"
	"basic_server/server/request"
	"basic_server/server/response"
	"basic_server/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

func (handler PostHandler) SavePost() gin.HandlerFunc {
	return func(context *gin.Context) {
		var createPostRequest request.CreatePostRequest
		if err := context.ShouldBind(&createPostRequest); err != nil {
			response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
			return
		}

		postService := service.PostService{}
		post := postService.CreatePost(createPostRequest.Title, createPostRequest.Content)
		handler.DB.Create(&post)
		response.SuccessResponse(context, response.CreatePostResponse{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
		})
	}
}

// GetPosts godoc
// @Summary Get all posts
// @Description Get all posts of all users
// @ID get-posts
// @Tags Post Actions
// @Produce json
// @Success 200 {object} response.CollectionResponse
// @Failure 401 {object} response.Error
// @Security ApiKeyAuth
// @Router /posts [get]
func (handler PostHandler) GetPosts() gin.HandlerFunc {
	return func(context *gin.Context) {
		postsRepository := repository.PostRepository{DB: handler.DB}
		var posts []model.Post
		postsRepository.GetAll(&posts)
		response.SuccessResponse(context, response.CreatePostsCollectionResponse(posts))
	}
}
