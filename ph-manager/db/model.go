package db

import "time"

type (
	Recording struct {
		ID                 int
		VideoName          string
		OriginalFileName   string
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
