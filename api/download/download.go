package download

import (
	"bigyohann/apidownloader/onefichier"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostDownload struct {
	Url string `json:"url"`
}

func DownloadHandler(c *gin.Context) {
	var json PostDownload
	c.BindJSON(&json)
	fileData, err := onefichier.GetFileData(json.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go onefichier.DownloadFile(fileData.Url, fileData)
	c.JSON(http.StatusOK, fileData)
}

func DataHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "data"})
}
