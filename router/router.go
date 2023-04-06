package router

import (
	"github.com/gallifreyCar/gin-file-server/handler"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	r := gin.Default()
	r.POST("/uploadFile", handler.UploadFileSingle)
	r.POST("/uploadFiles", handler.UploadFiles)
	r.GET("/downloadFile/:folder/:file_name", handler.DownloadFile)

	return r
}
