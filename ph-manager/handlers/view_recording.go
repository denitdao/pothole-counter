package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"strconv"
)

type (
	ViewRecordingPage struct {
		RecordingID int
		VideoName   string
		Detections  []Detection
		Error       error
	}

	Detection struct {
		ID               int
		FileName         string
		Confidence       float32
		FrameNumber      int
		TotalFrameNumber int
	}
)

func ViewRecording(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		renderFailureVR(c, err)
		return
	}

	recording, err := db.GetRecording(id)
	if err != nil {
		renderFailureVR(c, err)
		return
	}

	detections, err := db.GetDetections(id)
	if err != nil {
		renderFailureVR(c, err)
		return
	}

	viewDetections := make([]Detection, len(detections))
	for i, detection := range detections {
		viewDetections[i] = Detection{
			ID:               detection.ID,
			FileName:         detection.FileName,
			Confidence:       detection.Confidence,
			FrameNumber:      detection.FrameNumber,
			TotalFrameNumber: detection.TotalFrameNumber,
		}
	}

	p := ViewRecordingPage{
		RecordingID: id,
		VideoName:   recording.VideoName,
		Detections:  viewDetections,
	}
	c.HTML(http.StatusOK, "view-recording.gohtml", p)
}

func renderFailureVR(c *gin.Context, err error) {
	c.HTML(http.StatusBadRequest, "view-recording.gohtml", UploadRecordingComponent{
		Error: err,
	})
}
