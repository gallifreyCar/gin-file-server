package router

import (
	"github.com/gallifreyCar/gin-file-server/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

// MaxAllowed Middleware function to set a size limit on the uploaded files
func MaxAllowed(n int64) gin.HandlerFunc {
	maxBytes := n
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		if err := c.Request.ParseMultipartForm(maxBytes); err != nil {
			if err.Error() == "http: request body too large" {
				c.AbortWithStatus(http.StatusRequestEntityTooLarge)
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}

func setupRouter() *gin.Engine {

	r := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	//set a size limit on the uploaded files
	r.Use(MaxAllowed(9 << 20))
	r.POST("/uploadFile", handler.UploadFileSingle)
	r.POST("/uploadFiles", handler.UploadFiles)
	r.GET("/downloadFile/:folder/:file_name", handler.DownloadFile)

	return r
}
