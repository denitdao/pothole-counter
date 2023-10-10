package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type (
	UploadGpxComponent struct {
		RecordingID int
		Error       error
	}
)

func UploadGpx(c *gin.Context) {
	recordingID, err := strconv.Atoi(c.Param("recording_id"))
	if err != nil {
		renderFailureUG(c, err) // TODO: general 404 page
		return
	}

	c.Status(http.StatusOK)
	c.Header("HX-Redirect", fmt.Sprintf("/view-recording/%d", recordingID))
}

func renderFailureUG(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "upload-gpx.gohtml", UploadGpxComponent{
		Error: err,
	})
}
