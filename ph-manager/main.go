package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Record struct {
	Text  string
	Order int
}

func main() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		data := struct {
			Records []Record
		}{
			Records: []Record{
				{Text: "First", Order: 1},
				{Text: "Second", Order: 2},
			},
		}
		tmpl.Execute(w, data)
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		text := r.PostFormValue("text")
		order := r.PostFormValue("order")
		htmlStr := fmt.Sprintf("<div><h1>%s</h1><p>%s</p></div>", text, order)
		tmpl, _ := template.New("record").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/records", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
