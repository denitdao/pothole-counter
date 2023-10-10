package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddGpx(c *gin.Context) {
	recordingID, err := strconv.Atoi(c.Param("recording_id"))
	if err != nil {
		renderFailureUG(c, err)
		return
	}

	c.HTML(http.StatusOK, "add-gpx.gohtml", UploadGpxComponent{
		RecordingID: recordingID,
	})
}
