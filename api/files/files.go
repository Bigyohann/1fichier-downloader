package files

import (
	"bigyohann/apidownloader/internal/database"
	"bigyohann/apidownloader/internal/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DataHandler(c *gin.Context) {
	var files []models.File
	database.GetDB().Find(&files)
	c.JSON(http.StatusOK, files)
}
