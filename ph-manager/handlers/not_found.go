package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	ErrorPage struct {
		Error error
	}
)

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.gohtml", ErrorPage{})
}
