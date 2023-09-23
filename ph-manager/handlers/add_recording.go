package handlers

import (
	"html/template"
	"net/http"
)

func AddRecording(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/add-recording.gohtml", "templates/components/meta.gohtml"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
