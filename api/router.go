// Package api
package api

import (
	"bigyohann/apidownloader/api/download"
	"bigyohann/apidownloader/api/files"

	"github.com/gin-gonic/gin"
)

func HandleRouter(r *gin.Engine) *gin.Engine {
	r.POST("/download/:provider/get", download.DownloadHandler)
	r.POST("/download/:provider/data", download.DataHandler)

	r.GET("/files", files.GetAllFiles)
	r.GET("/files/downloading", files.GetDowloadingFiles)
	return r
}
