package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime"
	"net/http"
	"path"
)

func uploadFileSingle(c *gin.Context) {

	// Get the field name for file uploads from the request
	fieldName := c.DefaultPostForm("fieldName", "file")
	// Single file
	file, _ := c.FormFile(fieldName)
	log.Println(file.Filename)

	dst := "../target/single/" + file.Filename
	// Save the uploaded file to the specified directory
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		log.Fatal(err)
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func uploadFiles(c *gin.Context) {

	// Get the field name for file uploads from the request
	fieldName := c.DefaultPostForm("fieldName", "files")
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	// Get the uploaded files based on the specified field name
	files := form.File[fieldName]
	for _, file := range files {
		log.Println(file.Filename)

		dst := "../target/multiple/" + file.Filename
		// Save the uploaded file to the specified directory
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			log.Fatal(err)
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func downloadFile(c *gin.Context) {
	folder := c.Param("folder")
	fileName := c.Param("file_name")
	filePath := path.Join("..", "target", folder, fileName)
	ext := path.Ext(filePath)
	c.Header("Content-Type", mime.TypeByExtension(ext))
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Status(http.StatusOK)
	c.File(filePath)
}
