package dq

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Push 入队
func Push(c *gin.Context) {
	message := "Push"
	c.String(http.StatusOK, message)
}

// Pop 出队
func Pop(c *gin.Context) {
	message := "Push"
	c.String(http.StatusOK, message)
}
