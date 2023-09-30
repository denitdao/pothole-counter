package db

import (
	"time"
)

func CreateRecording(videoName string, originalFileName string, createdAt time.Time) (Recording, error) {
	result, err := DB.Exec(
		"INSERT INTO recordings (video_name, original_file_name, created_at) VALUES (?, ?, ?)",
		videoName, originalFileName, createdAt,
	)
	if err != nil {
		return Recording{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Recording{}, err
	}

	return Recording{
		ID:               int(id),
		VideoName:        videoName,
		OriginalFileName: originalFileName,
		CreatedAt:        createdAt,
	}, nil
}

func GetRecordings() ([]Recording, error) {
	rows, err := DB.Query("SELECT id, video_name, original_file_name, created_at FROM recordings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recordings []Recording
	for rows.Next() {
		var recording Recording
		err := rows.Scan(&recording.ID, &recording.VideoName, &recording.OriginalFileName, &recording.CreatedAt)
		if err != nil {
			return nil, err
		}
		recordings = append(recordings, recording)
	}

	return recordings, nil
}

func GetDetections(recordingID int) ([]Detection, error) {
	rows, err := DB.Query(`
		SELECT id, 
		       recording_id, 
		       file_name, 
		       frame_number, 
		       total_frame_number, 
		       video_millisecond, 
		       total_video_millisecond,
		       confidence, 
		       created_at 
		FROM detections 
		WHERE recording_id = ?
		`, recordingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var detections []Detection
	for rows.Next() {
		var detection Detection
		err := rows.Scan(&detection.ID,
			&detection.RecordingID,
			&detection.FileName,
			&detection.FrameNumber,
			&detection.TotalFrameNumber,
			&detection.VideoMillisecond,
			&detection.TotalVideoMillisecond,
			&detection.Confidence,
			&detection.CreatedAt)
		if err != nil {
			return nil, err
		}
		detections = append(detections, detection)
	}

	return detections, nil
}
