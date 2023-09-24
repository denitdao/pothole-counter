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
		Potholes int
		HasGPX   bool
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
			Potholes: -1,    // todo: get potholes
			HasGPX:   false, // todo: get gpx
			DateTime: recording.CreatedAt,
		}
	}

	c.HTML(http.StatusOK, "index.gohtml", p)
}
