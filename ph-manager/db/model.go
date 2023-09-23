package db

import "time"

type (
	Recording struct {
		ID               int
		VideoName        string
		OriginalFileName string
		CreatedAt        time.Time
	}

	Detection struct {
		ID               int
		RecordingID      int
		FileName         string
		FrameNumber      int
		VideoMillisecond int
		Confidence       float32
		CreatedAt        time.Time
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
		GpxID       int
		Latitude    float64
		Longitude   float64
		CreatedAt   time.Time
	}
)
