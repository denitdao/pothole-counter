package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"strconv"
)

type (
	ViewDetectionComponent struct {
		RecordingID int
		ID          int
		FileName    string
		Confidence  float32
		Latitude    float64
		Longitude   float64

		Error error
	}
)

func ViewDetection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		renderFailureVD(c, err)
		return
	}

	detection, err := db.GetDetection(id)
	if err != nil {
		renderFailureVD(c, err)
		return
	}

	location, err := db.GetLocationByDetectionID(id)
	if err != nil {
		renderFailureVD(c, err)
		return
	}

	c.HTML(http.StatusOK, "view-detection.gohtml", ViewDetectionComponent{
		RecordingID: detection.RecordingID,
		ID:          detection.ID,
		FileName:    detection.FileName,
		Confidence:  detection.Confidence,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
	})
}

func renderFailureVD(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "view-detection.gohtml", ViewDetectionComponent{
		Error: err,
	})
}
