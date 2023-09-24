package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func ViewRecording(c *gin.Context) {
	id := strings.TrimPrefix(c.Request.URL.Path, "/view-recording/")
	log.Println(id)

	c.HTML(http.StatusOK, "view-recording.gohtml", nil)
}
