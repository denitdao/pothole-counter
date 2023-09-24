package handlers

import (
	"html/template"
	"net/http"
	"ph-manager/util"
)

func AddRecording(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("add-recording.gohtml").Funcs(template.FuncMap{
		"formatDate": util.FormatDate,
	}).ParseFiles(
		"templates/add-recording.gohtml",
		"templates/components/meta.gohtml",
		"templates/components/upload-recording.gohtml",
	))
	err := tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
