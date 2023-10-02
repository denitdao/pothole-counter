package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"path/filepath"
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
	router.Static("/images",
		filepath.Join(util.GetProperty("storage.path"), util.GetProperty("record.folder")))
	router.Static("/videos",
		filepath.Join(util.GetProperty("storage.path"), util.GetProperty("video.folder")))
	router.SetFuncMap(template.FuncMap{
		"formatDate":  util.FormatDate,
		"mul":         util.Mul,
		"formatFloat": util.FormatFloat,
	})
	router.LoadHTMLGlob("templates/**/*.gohtml")

	// Setup routes
	router.GET("/", handlers.Index)
	router.GET("/add-recording", handlers.AddRecording)
	router.GET("/view-recording/:id", handlers.ViewRecording)
	router.POST("/upload-recording", handlers.UploadRecording)
	router.POST("/analyze/:id", handlers.AnalyzeRecording)
	router.DELETE("/detection/:id", handlers.DeleteDetection)

	// Start server
	log.Fatal(router.Run(":8080"))
}
