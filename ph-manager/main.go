package main

import (
	"log"
	"net/http"
	"ph-manager/handlers"
	"ph-manager/util"
)

func main() {
	// Initialize database
	util.InitDB()

	// Handle static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Setup routes
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/add-recording", handlers.AddRecording)
	http.HandleFunc("/view-recording/", handlers.ViewRecording)
	http.HandleFunc("/upload-recording", handlers.UploadRecording)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
