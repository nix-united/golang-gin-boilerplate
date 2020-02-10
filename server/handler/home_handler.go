package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeHandler struct {}

func (handler HomeHandler) Index() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(http.StatusOK, "Hello")
	}
}