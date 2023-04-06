package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func uploadFileSingle(c *gin.Context) {

	// Get the field name for file uploads from the request
	fieldName := c.DefaultPostForm("fieldName", "file")
	// Single file
	file, _ := c.FormFile(fieldName)
	log.Println(file.Filename)

	dst := "../target/single/" + file.Filename
	// Save the uploaded file to the specified directory
	c.SaveUploadedFile(file, dst)

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
		c.SaveUploadedFile(file, dst)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
