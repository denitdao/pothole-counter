package handlers

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ph-manager/db"
	"ph-manager/util"
	"time"
)

// TODO: steps to implement this page
//  make reusable component to inject into the page and return from endpoint

type (
	UploadRecordingComponent struct {
		UploadStatus UploadStatus
		Error        error
	}

	UploadStatus struct {
		RecordingID int
		VideoName   string
		Duration    int
		UploadedAt  time.Time

		HasVideo bool
		HasGPX   bool
	}
)

func UploadRecording(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	uploadStatus := UploadStatus{}

	recording, err := storeVideo(r)
	if err != nil {
		uploadStatus.HasVideo = false
	} else {
		uploadStatus.HasVideo = true
		uploadStatus.RecordingID = recording.ID
		uploadStatus.VideoName = recording.OriginalFileName
		uploadStatus.Duration = -1 // todo: get duration
		uploadStatus.UploadedAt = recording.CreatedAt
	}

	t := template.Must(template.New("upload-recording.gohtml").Funcs(template.FuncMap{
		"formatDate": util.FormatDate,
	}).ParseFiles("templates/components/upload-recording.gohtml"))
	c := UploadRecordingComponent{
		UploadStatus: uploadStatus,
		Error:        err,
	}
	err = t.Execute(w, c)

	if err != nil {
		panic(err)
	}
}

func storeVideo(r *http.Request) (db.Recording, error) {
	videoFile, header, err := r.FormFile("video")
	if err != nil {
		log.Println(err)
		return db.Recording{}, errors.New("unable to read video")
	}
	defer videoFile.Close()

	originalFileName := header.Filename

	storagePath, videoFolder := util.GetProperty("storage.path"), util.GetProperty("video.folder")
	videoDestPath := filepath.Join(storagePath, videoFolder)
	videoDest, err := os.CreateTemp(videoDestPath, "v_*.mp4")
	if err != nil {
		return db.Recording{}, errors.New("unable to save video")
	}
	defer videoDest.Close()
	_, err = io.Copy(videoDest, videoFile)
	if err != nil {
		return db.Recording{}, errors.New("unable to save video")
	}

	return db.CreateRecording(filepath.Base(videoDest.Name()), originalFileName, time.Now())
}

// isGPXPresent(w, r) && storeGPX(w, r)
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
