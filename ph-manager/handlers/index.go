package handlers

import (
	"html/template"
	"net/http"
	"time"
)

type (
	IndexPage struct {
		Recordings []Recording
	}

	Recording struct {
		ID       int
		Potholes int
		HasGPX   bool
		DateTime time.Time
	}
)

func Index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index.gohtml").Funcs(template.FuncMap{
		"formatDate": formatDate,
	}).ParseFiles("templates/index.gohtml", "templates/components/meta.gohtml"))

	p := IndexPage{
		Recordings: []Recording{
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

func formatDate(t time.Time) string {
	return t.Format("15:04 02.01.2006")
}
