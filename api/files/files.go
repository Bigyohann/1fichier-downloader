package files

import (
	"bigyohann/apidownloader/internal/database"
	"bigyohann/apidownloader/internal/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DataHandler(c *gin.Context) {
	var files []models.File
	// order by id desc
	database.GetDB().Order("created_at desc").Order("downloaded asc").Find(&files)
	c.JSON(http.StatusOK, files)
}
