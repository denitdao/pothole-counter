package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"strconv"
)

func AddGpx(c *gin.Context) {
	recordingID, err := strconv.Atoi(c.Param("recording_id"))
	if err != nil {
		NotFound(c)
		return
	}

	recording, err := db.GetRecording(recordingID)
	if err != nil || recording.Type != "VIDEO" {
		NotFound(c)
		return
	}

	c.HTML(http.StatusOK, "add-gpx.gohtml", UploadGpxComponent{
		RecordingID: recordingID,
	})
}
