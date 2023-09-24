package handlers

import (
	"html/template"
	"net/http"
	"ph-manager/db"
	"ph-manager/util"
	"time"
)

type (
	IndexPage struct {
		RecordingRows []RecordingRow
	}

	RecordingRow struct {
		ID       int
		Potholes int
		HasGPX   bool
		DateTime time.Time
	}
)

func Index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index.gohtml").Funcs(template.FuncMap{
		"formatDate": util.FormatDate,
	}).ParseFiles("templates/index.gohtml", "templates/components/meta.gohtml"))

	recordings, err := db.GetRecordings()
	if err != nil {
		panic(err)
	}

	p := IndexPage{
		RecordingRows: make([]RecordingRow, len(recordings)),
	}
	for i, recording := range recordings {
		p.RecordingRows[i] = RecordingRow{
			ID:       recording.ID,
			Potholes: -1,    // todo: get potholes
			HasGPX:   false, // todo: get gpx
			DateTime: recording.CreatedAt,
		}
	}

	err = t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}
