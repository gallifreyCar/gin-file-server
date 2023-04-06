package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func uploadFileSingle(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	dst := "../target/" + file.Filename
	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(file, dst)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func uploadFiles(c *gin.Context) {
	// Multipart form
	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	files := form.File["upload[]"]
	for _, file := range files {
		log.Println(file.Filename)

		dst := "../target/" + file.Filename
		// 上传文件至指定目录
		c.SaveUploadedFile(file, dst)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
