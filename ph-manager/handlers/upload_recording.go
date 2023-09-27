package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ph-manager/db"
	"ph-manager/util"
	"time"
)

type (
	UploadRecordingComponent struct {
		UploadStatus UploadStatus
		Error        error
	}

	UploadStatus struct {
		RecordingID int
		VideoName   string
		Duration    int
		Processing  bool
		UploadedAt  time.Time

		HasVideo bool
		HasGPX   bool
	}
)

func UploadRecording(c *gin.Context) {
	uploadStatus := UploadStatus{}

	recording, err := storeVideo(c.Request)
	if err != nil {
		uploadStatus.HasVideo = false
	} else {
		uploadStatus.HasVideo = true
		uploadStatus.RecordingID = recording.ID
		uploadStatus.VideoName = recording.OriginalFileName
		uploadStatus.Duration = -1 // todo: get duration
		uploadStatus.UploadedAt = recording.CreatedAt

		resp, err := http.Post(fmt.Sprintf("%s/analyze/%d", util.GetProperty("ph.detector.url"), recording.ID), "application/json", nil)
		if err == nil && resp.StatusCode == http.StatusOK {
			uploadStatus.Processing = true
			println("Successfully sent recording to analyzer")
		}
	}

	c.HTML(http.StatusOK, "upload-recording.gohtml", UploadRecordingComponent{
		UploadStatus: uploadStatus,
		Error:        err,
	})
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

// TODO: isGPXPresent(w, r) && storeGPX(w, r)
