package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"time"
)

type (
	IndexPage struct {
		RecordingRows []RecordingRow
	}

	RecordingRow struct {
		ID       int
		Status   string
		Potholes int
		DateTime time.Time
	}
)

func Index(c *gin.Context) {
	recordings, err := db.GetRecordings()
	if err != nil {
		panic(err)
	}

	p := IndexPage{
		RecordingRows: make([]RecordingRow, len(recordings)),
	}
	for i, recording := range recordings {
		p.RecordingRows[i] = RecordingRow{
			ID:       recording.ID,
			Status:   recording.Status,
			Potholes: recording.NumberOfDetections,
			DateTime: recording.CreatedAt,
		}
	}

	c.HTML(http.StatusOK, "index.gohtml", p)
}
