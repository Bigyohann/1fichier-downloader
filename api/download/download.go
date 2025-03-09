package download

import (
	"bigyohann/apidownloader/internal/service/downloader"
	"bigyohann/apidownloader/pkg/onefichier"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PostDownload struct {
	Url string `json:"url"`
}

func DownloadHandler(c *gin.Context) {
	var json PostDownload
	c.BindJSON(&json)
	json.Url = sanitizeUrl(json.Url)
	fileData, err := onefichier.GetFileData(json.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// exemple of url https://1fichier.com/?imlkq5gsl1igb7vu57kp&af=651048
	// we need to remove the &af=value

	file := downloader.HandleDownloadFile(fileData.Url)
	c.JSON(http.StatusOK, file)
}

func DataHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "data"})
}

func sanitizeUrl(url string) string {
	if strings.Contains(url, "&af=") {
		url = strings.Split(url, "&af=")[0]
	}
	return url
}
