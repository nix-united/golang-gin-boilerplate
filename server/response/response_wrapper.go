package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(context *gin.Context, statusCode int, data interface{}) {
	context.JSON(statusCode, data)
}

func SuccessResponse(context *gin.Context, data interface{}) {
	Response(context, http.StatusOK, data)
}

func ErrorResponse(context *gin.Context, statusCode int, message string) {
	Response(context, statusCode, struct {
		Code  int
		Error string
	}{Code: http.StatusBadRequest, Error: message})
}
