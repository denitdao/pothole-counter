package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"strconv"
	"strings"
)

type (
	ViewRecordingPage struct {
		RecordingID int
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
	id, err := strconv.Atoi(strings.TrimPrefix(c.Request.URL.Path, "/view-recording/"))

	if err != nil {
		c.HTML(http.StatusOK, "view-recording.gohtml", ViewRecordingPage{
			Error: err,
		})
		return
	}

	detections, err := db.GetDetections(id)
	if err != nil {
		c.HTML(http.StatusOK, "view-recording.gohtml", ViewRecordingPage{
			Error: err,
		})
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
		Detections:  viewDetections,
	}
	c.HTML(http.StatusOK, "view-recording.gohtml", p)
}
