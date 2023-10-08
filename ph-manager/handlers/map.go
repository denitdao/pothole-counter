package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
)

func ViewMap(c *gin.Context) {
	c.HTML(http.StatusOK, "view-map.gohtml", nil)
}

func GetMapJsonData(c *gin.Context) {
	locations, err := db.GetLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, locations)
}
