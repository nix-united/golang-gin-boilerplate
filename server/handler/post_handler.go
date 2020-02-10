package handler

import (
	"basic_server/server/model"
	"basic_server/server/repository"
	"basic_server/server/request"
	"basic_server/server/response"
	"basic_server/server/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
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
			Id:      post.ID,
			Title:   post.Title,
			Content: post.Content,
		})
	}
}

func (handler PostHandler) GetPosts() gin.HandlerFunc {
	return func(context *gin.Context) {
		postsRepository := repository.PostRepository{DB:handler.DB}
		var posts []model.Post
		postsRepository.GetAll(&posts)
		response.SuccessResponse(context, response.CreatePostsCollectionResponse(posts))
	}
}