package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"strconv"
)

type (
	ViewDetectionComponent struct {
		Detections []Detection

		Error error
	}
)

func ViewDetection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		renderFailureVD(c, err)
		return
	}

	location, err := db.GetLocationByDetectionID(id)
	if err != nil {
		renderFailureVD(c, err)
		return
	}

	allDetections, err := db.GetDetectionsAtLocation(location.Latitude, location.Longitude)
	if err != nil {
		renderFailureVD(c, err)
		return
	}

	var paramDetections []Detection
	for _, d := range allDetections {
		paramDetections = append(paramDetections, Detection{
			ID:               d.ID,
			RecordingID:      d.RecordingID,
			FileName:         d.FileName,
			Confidence:       d.Confidence,
			FrameNumber:      d.FrameNumber,
			TotalFrameNumber: d.TotalFrameNumber,
			Latitude:         location.Latitude,
			Longitude:        location.Longitude,
		})
	}

	c.HTML(http.StatusOK, "view-detection.gohtml", ViewDetectionComponent{
		Detections: paramDetections,
	})
}

func renderFailureVD(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "view-detection.gohtml", ViewDetectionComponent{
		Error: err,
	})
}
