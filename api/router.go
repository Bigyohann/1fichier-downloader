package api

import (
	"bigyohann/apidownloader/api/download"
	"bigyohann/apidownloader/api/files"

	"github.com/gin-gonic/gin"
)

// import routes from download.go
func HandleRouter(r *gin.Engine) *gin.Engine {
	r.POST("/download/get", download.DownloadHandler)
	r.GET("/download/data", download.DataHandler)

	r.GET("/files", files.DataHandler)
	return r
}
