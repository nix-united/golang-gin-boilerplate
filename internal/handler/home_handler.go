package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (handler HomeHandler) Index() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(http.StatusOK, "Hello")
	}
}
