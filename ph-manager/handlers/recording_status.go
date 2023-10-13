package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"ph-manager/util"
	"strconv"
)

type (
	RecordingStatusComponent struct {
		RecordingID int
		Type        string
		Status      string
		GpxStatus   string
	}
)

func AnalyzeRecording(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		renderFailureRS(c)
		return
	}

	recording, err := db.GetRecording(id)
	if err != nil {
		renderFailureRS(c)
		return
	}

	gpxStatus := "MISSING"
	gpx, err := db.GetGpxID(id)
	if err == nil {
		gpxStatus = gpx.Status
	}

	resp, err := http.Post(fmt.Sprintf("%s/analyze/%d", util.GetProperty("ph.detector.url"), id), "application/json", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		println("Failed to send recording to analyzer")
		c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
			RecordingID: id,
			Type:        recording.Type,
			Status:      "FAILED",
			GpxStatus:   gpxStatus,
		})
		return
	}

	c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
		RecordingID: id,
		Type:        recording.Type,
		Status:      "PROCESSING",
		GpxStatus:   gpxStatus,
	})
}

func LocateRecording(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		renderFailureRS(c)
		return
	}

	recording, err := db.GetRecording(id)
	if err != nil {
		renderFailureRS(c)
		return
	}

	if recording.Type != "VIDEO" {
		c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
			RecordingID: id,
			Type:        recording.Type,
			Status:      recording.Status,
			GpxStatus:   "MISSING",
		})
		return
	}

	_, err = db.GetGpxID(id)
	if err != nil {
		c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
			RecordingID: id,
			Type:        recording.Type,
			Status:      recording.Status,
			GpxStatus:   "MISSING",
		})
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/locate/%d", util.GetProperty("ph.detector.url"), id), "application/json", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		println("Failed to send recording to locator")
		c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
			RecordingID: id,
			Type:        recording.Type,
			Status:      recording.Status,
			GpxStatus:   "FAILED",
		})
		return
	}

	c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
		RecordingID: id,
		Type:        recording.Type,
		Status:      recording.Status,
		GpxStatus:   "PROCESSING",
	})
}

func renderFailureRS(c *gin.Context) {
	c.HTML(http.StatusBadRequest, "recording-status.gohtml", nil)
}
