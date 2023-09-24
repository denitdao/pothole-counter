package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"ph-manager/db"
	"ph-manager/handlers"
	"ph-manager/util"
)

func main() {
	// Initialize database
	db.InitDB()

	// Handle static files
	router := gin.Default()
	router.Static("/static", "./static")
	router.SetFuncMap(template.FuncMap{
		"formatDate": util.FormatDate,
	})
	router.LoadHTMLGlob("templates/**/*.gohtml")

	// Setup routes
	router.GET("/", handlers.Index)
	router.GET("/add-recording", handlers.AddRecording)
	router.GET("/view-recording/:id", handlers.ViewRecording)
	router.POST("/upload-recording", handlers.UploadRecording)

	// Start server
	log.Fatal(router.Run(":8080"))
}
