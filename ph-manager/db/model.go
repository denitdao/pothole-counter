package db

import (
	"database/sql"
	"time"
)

type (
	Recording struct {
		ID                 int
		FileName           string
		OriginalFileName   string
		Note               string
		Type               string
		Status             string
		CreatedAt          time.Time
		NumberOfDetections int
	}

	Detection struct {
		ID                    int
		RecordingID           int
		FileName              string
		FrameNumber           int
		TotalFrameNumber      int
		VideoMillisecond      int
		TotalVideoMillisecond int
		Confidence            float32
		CreatedAt             time.Time

		DetectionLocation DetectionLocation
	}

	GPX struct {
		ID          int
		RecordingID int
		FileName    string
		CreatedAt   time.Time
	}

	DetectionLocation struct {
		ID          int
		DetectionID int
		GpxID       sql.NullInt64
		Latitude    float64
		Longitude   float64
		CreatedAt   time.Time
	}
)
