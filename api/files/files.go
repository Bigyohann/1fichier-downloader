// Package files handles file-related operations in the API.
package files

import (
	"net/http"

	"bigyohann/apidownloader/internal/database"
	"bigyohann/apidownloader/internal/database/models"

	"github.com/gin-gonic/gin"
)

func GetAllFiles(c *gin.Context) {
	var files []models.File
	// order by id desc
	database.GetDB().Order("created_at desc").Order("downloaded asc").Find(&files)
	c.JSON(http.StatusOK, files)
}

func GetDowloadingFiles(c *gin.Context) {
	var files []models.File
	// order by id desc
	database.GetDB().
		Where("downloaded = ?", false).
		Where("status = ?", "Downloading").
		Order("created_at desc").
		Find(&files)
	c.JSON(http.StatusOK, files)
}
