package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ph-manager/util"
)

func UploadRecording(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	row := util.DB.QueryRow("SELECT id, video_name FROM recordings WHERE id = ?", "1000000")
	var id, videoName string
	err := row.Scan(&id, &videoName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id, videoName)

	if storeVideo(w, r) {
		return
	}
	if isGPXPresent(w, r) && storeGPX(w, r) {
		return
	}

	fmt.Fprint(w, "Upload successful")
	return
}

func storeVideo(w http.ResponseWriter, r *http.Request) bool {
	storagePath := util.GetProperty("storage.path")
	videoFolder := util.GetProperty("video.folder")
	videoDestPath := filepath.Join(storagePath, videoFolder, "video.mp4") // todo: unique filename

	videoFile, _, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Unable to read video", http.StatusBadRequest)
		return true
	}
	defer videoFile.Close()

	videoDest, err := os.Create(videoDestPath)
	if err != nil {
		http.Error(w, "Unable to save video", http.StatusInternalServerError)
		return true
	}
	defer videoDest.Close()
	io.Copy(videoDest, videoFile)
	// todo: create recording in database
	return false
}

func storeGPX(w http.ResponseWriter, r *http.Request) bool {
	storagePath := util.GetProperty("storage.path")
	gpxFolder := util.GetProperty("gpx.folder")
	gpxDestPath := filepath.Join(storagePath, gpxFolder, "file.gpx") // todo: unique filename

	gpxFile, _, err := r.FormFile("gpx")
	if err != nil {
		http.Error(w, "Unable to read GPX file", http.StatusBadRequest)
		return true
	}
	defer gpxFile.Close()

	gpxDest, err := os.Create(gpxDestPath)
	if err != nil {
		http.Error(w, "Unable to save GPX file", http.StatusInternalServerError)
		return true
	}
	defer gpxDest.Close()
	io.Copy(gpxDest, gpxFile)
	// todo: create gpx in database and link to recording
	return false
}

func isGPXPresent(w http.ResponseWriter, r *http.Request) bool {
	gpxFile, _, err := r.FormFile("gpx")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return false
		} else {
			http.Error(w, "Unable to read GPX file", http.StatusBadRequest)
			return false
		}
	}

	if gpxFile != nil {
		defer gpxFile.Close()
	}
	return true
}
