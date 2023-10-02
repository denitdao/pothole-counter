package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ph-manager/db"
	"strconv"
)

func DeleteDetection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = db.DeleteDetection(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
