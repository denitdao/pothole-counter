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
	"strconv"
	"time"
)

type (
	UploadGpxComponent struct {
		RecordingID int
		Success     bool
		Error       error
	}
)

func UploadGpx(c *gin.Context) {
	prop := UploadGpxComponent{}
	recordingID, err := strconv.Atoi(c.Param("recording_id"))
	if err != nil {
		prop.Error = err
		renderFailureUG(c, prop)
		return
	}
	prop.RecordingID = recordingID

	recording, err := db.GetRecording(recordingID)
	if err != nil || recording.Type != "VIDEO" {
		prop.Error = err
		renderFailureUG(c, prop)
		return
	}

	gpx, err := storeGpx(c.Request)
	if err != nil {
		prop.Error = err
		renderFailureUG(c, prop)
		return
	}
	prop.Success = true

	gpx.ID = recordingID
	gpx, err = db.CreateUpdateGpx(gpx)
	if err != nil {
		prop.Error = err
		renderFailureUG(c, prop)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/locate/%d", util.GetProperty("ph.detector.url"), recording.ID), "application/json", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		err = errors.New("failed to send gpx for locating")
		prop.Error = err
		renderFailureUG(c, prop)
		return
	}

	c.Status(http.StatusOK)
	c.Header("HX-Redirect", fmt.Sprintf("/view-recording/%d", recordingID))
}

func storeGpx(r *http.Request) (db.GPX, error) {
	gpxFile, _, err := r.FormFile("file")
	if err != nil {
		return db.GPX{}, errors.New("unable to read file")
	}
	defer gpxFile.Close()

	storagePath, gpxFolder := util.GetProperty("storage.path"), util.GetProperty("gpx.folder")
	gpxDestPath := filepath.Join(storagePath, gpxFolder)
	gpxDest, err := os.CreateTemp(gpxDestPath, "g_*.gpx")
	if err != nil {
		return db.GPX{}, errors.New("unable to save file")
	}
	defer gpxDest.Close()
	_, err = io.Copy(gpxDest, gpxFile)
	if err != nil {
		return db.GPX{}, errors.New("unable to save file")
	}

	return db.GPX{
		FileName:  filepath.Base(gpxDest.Name()),
		Status:    "CREATED",
		CreatedAt: time.Now(),
	}, nil
}

func renderFailureUG(c *gin.Context, prop UploadGpxComponent) {
	c.HTML(http.StatusOK, "upload-gpx.gohtml", prop)
}
