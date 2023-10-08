package db

func CreateRecording(recording Recording) (Recording, error) {
	result, err := DB.Exec(
		"INSERT INTO recordings (file_name, original_file_name, note, type, created_at) VALUES (?, ?, ?, ?, ?)",
		recording.FileName, recording.OriginalFileName, recording.Note, recording.Type, recording.CreatedAt,
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
	err := DB.QueryRow(`SELECT id, file_name, original_file_name, note, type, status, created_at FROM recordings WHERE id = ? and deleted = FALSE`, id).
		Scan(&recording.ID, &recording.FileName, &recording.OriginalFileName, &recording.Note, &recording.Type, &recording.Status, &recording.CreatedAt)
	if err != nil {
		return Recording{}, err
	}

	return recording, nil
}

func GetRecordings() ([]Recording, error) {
	rows, err := DB.Query(`
		SELECT r.id,
			   r.file_name,
			   original_file_name,
			   note,
			   type,
			   status,
			   r.created_at,
			   COUNT(d.id) AS number_of_detections
		FROM recordings r
			LEFT JOIN detections d ON r.id = d.recording_id AND d.deleted = FALSE
		WHERE r.deleted = FALSE
		GROUP BY r.id,
				 r.file_name,
				 original_file_name,
				 note,
				 type,
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
		err := rows.Scan(&recording.ID, &recording.FileName, &recording.OriginalFileName, &recording.Note, &recording.Type, &recording.Status, &recording.CreatedAt, &recording.NumberOfDetections)
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

func GetDetection(id int) (Detection, error) {
	var detection Detection
	err := DB.QueryRow(`
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
		WHERE id = ? and deleted = FALSE
		`, id).
		Scan(&detection.ID,
			&detection.RecordingID,
			&detection.FileName,
			&detection.FrameNumber,
			&detection.TotalFrameNumber,
			&detection.VideoMillisecond,
			&detection.TotalVideoMillisecond,
			&detection.Confidence,
			&detection.CreatedAt)
	if err != nil {
		return Detection{}, err
	}

	return detection, nil
}

func GetLocations() ([]DetectionLocation, error) {
	var locations []DetectionLocation
	rows, err := DB.Query(`
		SELECT l.id,
		       l.detection_id,
		       l.gpx_id,
		       l.latitude,
		       l.longitude,
		       l.created_at
		FROM detection_location l
			INNER JOIN detections d ON l.detection_id = d.id
		WHERE d.deleted = FALSE
		ORDER BY l.created_at
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var location DetectionLocation
		err := rows.Scan(&location.ID,
			&location.DetectionID,
			&location.GpxID,
			&location.Latitude,
			&location.Longitude,
			&location.CreatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func GetLocationByDetectionID(detectionID int) (DetectionLocation, error) {
	var location DetectionLocation
	err := DB.QueryRow(`
		SELECT id,
			   detection_id,
			   gpx_id,
			   latitude,
			   longitude,
			   created_at
		FROM detection_location
		WHERE detection_id = ?
		ORDER BY created_at
		`, detectionID).
		Scan(&location.ID,
			&location.DetectionID,
			&location.GpxID,
			&location.Latitude,
			&location.Longitude,
			&location.CreatedAt)
	if err != nil {
		return DetectionLocation{}, err
	}

	return location, nil
}

func GetLocationsByRecordingID(recordingID int) ([]DetectionLocation, error) {
	rows, err := DB.Query(`
        SELECT l.id,
               l.detection_id,
               l.gpx_id,
               l.latitude,
               l.longitude,
               l.created_at
        FROM detection_location l
        	INNER JOIN detections d ON l.detection_id = d.id
        WHERE d.recording_id = ? AND d.deleted = FALSE
        ORDER BY l.created_at
        `, recordingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []DetectionLocation
	for rows.Next() {
		var location DetectionLocation
		err := rows.Scan(&location.ID,
			&location.DetectionID,
			&location.GpxID,
			&location.Latitude,
			&location.Longitude,
			&location.CreatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func DeleteDetection(id int) error {
	_, err := DB.Exec("UPDATE detections SET deleted = TRUE WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
