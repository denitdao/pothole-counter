package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
		RecordingID       int
		OriginalVideoName string
	}
)

func UploadRecording(c *gin.Context) {
	uploadStatus := UploadStatus{}

	// TODO: check media type and apply correct strategy
	recording, err := storeVideo(c.Request)
	if err != nil {
		renderFailureUR(c, uploadStatus, err)
		return
	}
	uploadStatus.OriginalVideoName = recording.OriginalFileName

	recording, err = db.CreateRecording(recording)
	if err != nil {
		renderFailureUR(c, uploadStatus, err)
		return
	}
	uploadStatus.RecordingID = recording.ID

	resp, err := http.Post(fmt.Sprintf("%s/analyze/%d", util.GetProperty("ph.detector.url"), recording.ID), "application/json", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		err = errors.New("failed to send recording to analyzer")
		renderFailureUR(c, uploadStatus, err)
		return
	}

	c.Status(http.StatusOK)
	// TODO: redirect to /view-recording/video/:id or /view-recording/image/:id
	c.Header("HX-Redirect", fmt.Sprintf("/view-recording/%d", recording.ID))
}

func storeVideo(r *http.Request) (db.Recording, error) {
	videoFile, header, err := r.FormFile("video")
	if err != nil {
		return db.Recording{}, errors.New("unable to read video")
	}
	defer videoFile.Close()

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

	return db.Recording{
		FileName:         filepath.Base(videoDest.Name()),
		OriginalFileName: header.Filename,
		CreatedAt:        time.Now(),
		Note:             r.FormValue("note"),
		Type:             "VIDEO",
	}, nil
}

func renderFailureUR(c *gin.Context, uploadStatus UploadStatus, err error) {
	c.HTML(http.StatusOK, "upload-recording.gohtml", UploadRecordingComponent{
		UploadStatus: uploadStatus,
		Error:        err,
	})
}
