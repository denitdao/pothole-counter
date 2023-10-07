package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"strconv"
)

type (
	ViewRecordingPage struct {
		RecordingID      int
		FileName         string
		Note             string
		Status           string
		Detections       []Detection
		Error            error
		DetectionBatches []DetectionBatch
	}

	Detection struct {
		ID               int
		FileName         string
		Confidence       float32
		FrameNumber      int
		TotalFrameNumber int
		Latitude         float64
		Longitude        float64
	}

	DetectionBatch struct {
		StartFrame         int
		EndFrame           int
		NumberOfDetections int
		MaxBatchSize       int
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

	locations, err := db.GetLocationsByRecordingID(id)
	if err != nil {
		renderFailureVR(c, err)
		return
	}

	for i := range detections {
		for _, location := range locations {
			if detections[i].ID == location.DetectionID {
				detections[i].DetectionLocation = location
				break
			}
		}
	}

	viewDetections := make([]Detection, len(detections))
	for i, detection := range detections {
		viewDetections[i] = Detection{
			ID:               detection.ID,
			FileName:         detection.FileName,
			Confidence:       detection.Confidence,
			FrameNumber:      detection.FrameNumber,
			TotalFrameNumber: detection.TotalFrameNumber,
			Latitude:         detection.DetectionLocation.Latitude,
			Longitude:        detection.DetectionLocation.Longitude,
		}
	}

	p := ViewRecordingPage{
		RecordingID:      id,
		FileName:         recording.FileName,
		Note:             recording.Note,
		Status:           recording.Status,
		Detections:       viewDetections,
		DetectionBatches: calculateDetectionBatches(viewDetections),
	}

	switch recording.Type {
	case "VIDEO":
		c.HTML(http.StatusOK, "view-recording-video.gohtml", p)
	case "IMAGE":
		c.HTML(http.StatusOK, "view-recording-image.gohtml", p)
	}
}

func calculateDetectionBatches(detections []Detection) []DetectionBatch {
	if len(detections) == 0 {
		return nil
	}

	rangeSize := detections[0].TotalFrameNumber / 20
	batches := make([]DetectionBatch, 20)

	for i := 0; i < 20; i++ {
		batches[i].StartFrame = i * rangeSize
		batches[i].EndFrame = (i + 1) * rangeSize
	}

	for _, detection := range detections {
		for i, batch := range batches {
			if detection.FrameNumber >= batch.StartFrame && detection.FrameNumber <= batch.EndFrame {
				batches[i].NumberOfDetections++
			}
		}
	}

	maxBatchSize := 0
	for _, batch := range batches {
		if batch.NumberOfDetections > maxBatchSize {
			maxBatchSize = batch.NumberOfDetections
		}
	}

	for i := range batches {
		batches[i].MaxBatchSize = maxBatchSize
	}

	return batches
}

func renderFailureVR(c *gin.Context, err error) {
	c.HTML(http.StatusBadRequest, "view-recording-video.gohtml", UploadRecordingComponent{
		Error: err,
	})
}
