package db

func CreateRecording(recording Recording) (Recording, error) {
	result, err := DB.Exec(
		"INSERT INTO recordings (video_name, original_file_name, created_at) VALUES (?, ?, ?)",
		recording.VideoName, recording.OriginalFileName, recording.CreatedAt,
	)
	if err != nil {
		return Recording{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Recording{}, err
	}

	recording.ID = int(id)
	return recording, nil
}

func GetRecording(id int) (Recording, error) {
	var recording Recording
	err := DB.QueryRow(`SELECT id, video_name, original_file_name, status, created_at FROM recordings WHERE id = ? and deleted = FALSE`, id).
		Scan(&recording.ID, &recording.VideoName, &recording.OriginalFileName, &recording.Status, &recording.CreatedAt)
	if err != nil {
		return Recording{}, err
	}

	return recording, nil
}

func GetRecordings() ([]Recording, error) {
	rows, err := DB.Query(`
		SELECT r.id,
			   video_name,
			   original_file_name,
			   status,
			   r.created_at,
			   COUNT(d.id) AS number_of_detections
		FROM recordings r
			LEFT JOIN detections d ON r.id = d.recording_id AND d.deleted = FALSE
		WHERE r.deleted = FALSE
		GROUP BY r.id,
				 video_name,
				 original_file_name,
				 status,
				 r.created_at
		ORDER BY r.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recordings []Recording
	for rows.Next() {
		var recording Recording
		err := rows.Scan(&recording.ID, &recording.VideoName, &recording.OriginalFileName, &recording.Status, &recording.CreatedAt, &recording.NumberOfDetections)
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
		WHERE recording_id = ? and deleted = FALSE
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
