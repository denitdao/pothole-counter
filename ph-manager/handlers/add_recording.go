package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddRecording(c *gin.Context) {
	c.HTML(http.StatusOK, "add-recording.gohtml", nil)
}
