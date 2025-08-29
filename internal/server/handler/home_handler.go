package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h HomeHandler) Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}
