package handler

import (
	"basic_server/server/model"
	"basic_server/server/request"
	"basic_server/server/response"
	"basic_server/server/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type PostHandler struct {
	DB *gorm.DB
	PostService service.PostService
}

func (handler PostHandler) GetPostById() gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Param("id")
		post := handler.PostService.GetById(id)
		response.SuccessResponse(context, response.GetPostResponse{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
		})
	}
}

func (handler PostHandler) SavePost() gin.HandlerFunc {
	return func(context *gin.Context) {
		var createPostRequest request.CreatePostRequest
		if err := context.ShouldBind(&createPostRequest); err != nil {
			response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
			return
		}

		newPost := handler.PostService.CreatePost(createPostRequest.Title, createPostRequest.Content)
		post := handler.PostService.Save(newPost)
		response.SuccessResponse(context, response.CreatePostResponse{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
		})
	}
}

func (handler PostHandler) UpdatePost() gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Param("id")
		post := handler.PostService.GetById(id)
		if post == (model.Post{}) {
			context.Status(http.StatusBadRequest)
			return
		}

		post.Title = context.Param("title")
		post.Content = context.Param("content")
		handler.PostService.Save(post)
		response.SuccessResponse(context, response.GetPostResponse{
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
		posts := handler.PostService.GetAll()
		response.SuccessResponse(context, response.CreatePostsCollectionResponse(posts))
	}
}

func (handler PostHandler) DeletePost() gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Param("id")
		post := handler.PostService.GetById(id)
		if post == (model.Post{}) {
			context.Status(http.StatusBadRequest)
			return
		}

		handler.PostService.Delete(post)
		context.Status(http.StatusOK)
	}
}
