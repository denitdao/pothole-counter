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
		RecordingID:      id,
		FileName:         recording.FileName,
		Note:             recording.Note,
		Status:           recording.Status,
		Detections:       viewDetections,
		DetectionBatches: calculateDetectionBatches(viewDetections),
	}
	c.HTML(http.StatusOK, "view-recording.gohtml", p) // TODO: show correct slightly different template for video or image
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
	c.HTML(http.StatusBadRequest, "view-recording.gohtml", UploadRecordingComponent{
		Error: err,
	})
}
