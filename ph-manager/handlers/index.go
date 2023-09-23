package handlers

import (
	"html/template"
	"net/http"
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

	p := IndexPage{
		RecordingRows: []RecordingRow{
			{1, 3, true, time.Now()},
			{2, 5, false, time.Now()},
			{3, 1, false, time.Now()},
		},
	}

	err := t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}
