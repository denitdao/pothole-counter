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

	contentType, err := detectContentType(c)
	if err != nil {
		renderFailureUR(c, uploadStatus, err)
		return
	}
	var recording db.Recording
	switch {
	case contentType == "video/mp4":
		recording, err = storeVideo(c.Request)
		if err != nil {
			renderFailureUR(c, uploadStatus, err)
			return
		}
	case contentType == "image/jpeg":
		recording, err = storeImage(c.Request)
		if err != nil {
			renderFailureUR(c, uploadStatus, err)
			return
		}
	default:
		renderFailureUR(c, uploadStatus, fmt.Errorf("unsupported file type"))
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
	c.Header("HX-Redirect", fmt.Sprintf("/view-recording/%d", recording.ID))
}

func detectContentType(c *gin.Context) (string, error) {
	videoFile, _, err := c.Request.FormFile("file")
	if err != nil {
		return "", errors.New("unable to read video")
	}
	defer videoFile.Close()

	// Read the first 512 bytes to determine the file type
	buffer := make([]byte, 512)
	_, err = videoFile.Read(buffer)
	if err != nil {
		return "", errors.New("unable to read video")
	}
	// Reset the read pointer of the file to the beginning
	videoFile.Seek(0, 0)

	// Determine the content type
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

func storeVideo(r *http.Request) (db.Recording, error) {
	videoFile, header, err := r.FormFile("file")
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

func storeImage(r *http.Request) (db.Recording, error) {
	imageFile, header, err := r.FormFile("file")
	if err != nil {
		return db.Recording{}, errors.New("unable to read video")
	}
	defer imageFile.Close()

	storagePath, imageFolder := util.GetProperty("storage.path"), util.GetProperty("video.folder")
	imageDestPath := filepath.Join(storagePath, imageFolder)
	imageDest, err := os.CreateTemp(imageDestPath, "i_*.jpg")
	if err != nil {
		return db.Recording{}, errors.New("unable to save image")
	}
	defer imageDest.Close()
	_, err = io.Copy(imageDest, imageFile)
	if err != nil {
		return db.Recording{}, errors.New("unable to save image")
	}

	return db.Recording{
		FileName:         filepath.Base(imageDest.Name()),
		OriginalFileName: header.Filename,
		CreatedAt:        time.Now(),
		Note:             r.FormValue("note"),
		Type:             "IMAGE",
	}, nil
}

func renderFailureUR(c *gin.Context, uploadStatus UploadStatus, err error) {
	c.HTML(http.StatusOK, "upload-recording.gohtml", UploadRecordingComponent{
		UploadStatus: uploadStatus,
		Error:        err,
	})
}
