package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/util"
	"strconv"
)

type (
	RecordingStatusComponent struct {
		RecordingID int
		Status      string
	}
)

func AnalyzeRecording(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		renderFailureRS(c)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/analyze/%d", util.GetProperty("ph.detector.url"), id), "application/json", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		println("Failed to send recording to analyzer")
		c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
			RecordingID: id,
			Status:      "FAILED",
		})
		return
	}

	c.HTML(http.StatusOK, "recording-status.gohtml", RecordingStatusComponent{
		RecordingID: id,
		Status:      "PROCESSING",
	})
}

func renderFailureRS(c *gin.Context) {
	c.HTML(http.StatusBadRequest, "recording-status.gohtml", nil)
}
