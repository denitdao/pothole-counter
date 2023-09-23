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
