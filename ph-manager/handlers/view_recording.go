package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func ViewRecording(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/view-recording/")
	log.Println(id)
	tmpl := template.Must(template.ParseFiles("templates/view-recording.gohtml", "templates/components/meta.gohtml"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
