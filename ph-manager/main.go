package main

import (
	"log"
	"net/http"
	"ph-manager/db"
	"ph-manager/handlers"
)

func main() {
	// Initialize database
	db.InitDB()

	// Handle static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Setup routes
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/add-recording", handlers.AddRecording)
	http.HandleFunc("/view-recording/", handlers.ViewRecording)
	http.HandleFunc("/upload-recording", handlers.UploadRecording)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
