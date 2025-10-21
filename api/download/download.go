// Package download provides handlers for downloading files and retrieving file data.
package download

import (
	"net/http"
	"strings"

	"bigyohann/apidownloader/internal/service/downloader"

	"github.com/gin-gonic/gin"
)

type PostDownload struct {
	URL string `json:"url"`
}

func DownloadHandler(c *gin.Context) {
	downloadClient, err := getClientWithProvider(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	var json PostDownload
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	json.URL = sanitizeURL(json.URL)

	file := downloader.HandleDownloadFile(downloadClient, json.URL)

	c.JSON(http.StatusOK, file)
}

func getClientWithProvider(c *gin.Context) (downloader.DownloadClient, error) {
	downloadClient, err := GetClient(c.Param("provider"))
	if err != nil {
		return nil, err
	}
	return downloadClient, nil
}

func DataHandler(c *gin.Context) {
	downloadClient, err := getClientWithProvider(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	var json PostDownload
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fileData, _ := downloadClient.GetFileData(sanitizeURL(json.URL))

	c.JSON(http.StatusOK, fileData)
}

func sanitizeURL(url string) string {
	if strings.Contains(url, "&af=") {
		url = strings.Split(url, "&af=")[0]
	}
	return url
}
