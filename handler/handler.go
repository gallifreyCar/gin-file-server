package handler

import (
	"fmt"
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/gin-gonic/gin"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

func UploadFileSingle(c *gin.Context) {

	//set handle log
	logFile := log2.InitLogFile("handle.log", "[UploadFileSingle]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(logFile)

	// Get the field name for file uploads from the request
	fieldName := c.DefaultPostForm("fieldName", "file")

	// Single file
	file, err := c.FormFile(fieldName)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	dst := "../target/upload/single/" + file.Filename
	// Save the uploaded file to the specified directory
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	log.Println(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func UploadFiles(c *gin.Context) {

	//set handle log
	logFile := log2.InitLogFile("handle.log", "[UploadFileMultiple]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(logFile)

	// Get the field name for file uploads from the request
	fieldName := c.DefaultPostForm("fieldName", "files")
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))

		return
	}
	// Get the uploaded files based on the specified field name
	files := form.File[fieldName]
	for _, file := range files {
		log.Println(file.Filename)
		dst := "../target/upload/multiple/" + file.Filename
		// Save the uploaded file to the specified directory
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			log.Println(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))

			return
		}
	}
	log.Println(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func DownloadFile(c *gin.Context) {
	//set handle log
	logFile := log2.InitLogFile("handle.log", "[DownloadFile]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(logFile)

	// Get url param
	folder := c.Param("folder")
	fileName := c.Param("file_name")
	baseUrl := path.Join("..", "target", "upload")
	// Build local filePath
	filePath := path.Join(baseUrl, folder, fileName)
	// Check the files is existence or not
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println(err)
		// if file is not existence , return 404
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "file not found"})
		return
	}
	// Get ext
	ext := path.Ext(filePath)
	// Set response Header
	c.Header("Content-Type", mime.TypeByExtension(ext))
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Status(http.StatusOK)
	c.File(filePath)
	log.Printf("Folder: %v,Folder:%v,Code:%v\n", folder, folder, http.StatusOK)

}
