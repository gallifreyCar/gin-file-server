package handler

import (
	"fmt"
	"github.com/gallifreyCar/gin-file-server/m-logger"
	"github.com/gallifreyCar/gin-file-server/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime"
	"net/http"
	"os"
	"path"
)

func UploadFileSingle(c *gin.Context) {

	//set handle zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server.log", "[UploadFileSingle]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	// Get the field name for file uploads from the request
	fieldName := c.DefaultPostForm("fieldName", "file")

	// Single file
	file, err := c.FormFile(fieldName)
	if err != nil {
		logger.Error("Fail to get file from FormFile by fieldName", zap.Error(err))
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	dst := "../target/upload/single/" + file.Filename
	// Save the uploaded file to the specified directory
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		logger.Error("Fail to save upload file", zap.Int("statusCode", http.StatusBadRequest), zap.Error(err))
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	//insert an upload file zLogger record to database
	db, _ := repository.GetDataBase()
	userAgent := c.GetHeader("User-Agent")
	fileType := path.Ext(file.Filename)
	_, _, err = repository.InsertFileLog("../target/upload/single/", file.Filename, userAgent, fileType, file.Size, db)
	if err != nil {
		logger.Error("get form err: ", zap.Error(err))
	}
	logger.Info("File uploaded success", zap.Int("statusCode", http.StatusOK), zap.String("filename", file.Filename))
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func UploadFiles(c *gin.Context) {

	//set handle zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server.log", "[UploadFileMultiple]")
	defer closeFunc()

	// Get the field name for file uploads from the request
	fieldName := c.DefaultPostForm("fieldName", "files")
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error("Fail to get file from FormFile by fieldName", zap.Error(err))
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))

		return
	}
	// Get the uploaded files based on the specified field name
	files := form.File[fieldName]
	db, _ := repository.GetDataBase()
	for _, file := range files {
		logger.Info(file.Filename)
		dst := "../target/upload/multiple/" + file.Filename
		// Save the uploaded file to the specified directory
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			logger.Error("Fail to save upload file", zap.Int("statusCode", http.StatusBadRequest), zap.Error(err))
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))

			return
		}

		userAgent := c.GetHeader("User-Agent")
		fileType := path.Ext(file.Filename)
		_, _, err = repository.InsertFileLog("../target/upload/multiple/", file.Filename, userAgent, fileType, file.Size, db)
		if err != nil {
			logger.Error("Fail to insert file log into repository ", zap.Error(err))
		}

	}
	logger.Info("Upload file success!", zap.Int("statusCode", http.StatusOK))
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func DownloadFile(c *gin.Context) {

	//set handle zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server.log", "[DownloadFile]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	// Get url param
	folder := c.Param("folder")
	fileName := c.Param("file_name")
	baseUrl := path.Join("..", "target", "upload")
	// Build local filePath
	filePath := path.Join(baseUrl, folder, fileName)
	// Check the files is existence or not
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Error("get form err: ", zap.Error(err))
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
	logger.Info("Folder upload status",
		zap.String("folderName", folder),
		zap.String("folderPath", folder),
		zap.Int("statusCode", http.StatusOK))

}

func SelectFileLogByName(c *gin.Context) {
	//set handle  zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server.log", "[SelectFileLogByName]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()
	// Get url param
	fileName := c.Param("file_name")
	// Get logger
	db, err := repository.GetDataBase()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("get form err: %s", err.Error()))
		logger.Error("Fail to get database", zap.Error(err))

	}
	uploadFileLog, err := repository.SelectFileLog(fileName, db)
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("get form err: %s", err.Error()))
		logger.Error("Fail to select file log from repository", zap.Error(err))
	}

	c.IndentedJSON(http.StatusOK, uploadFileLog)
	logger.Info("Select file log success!", zap.Int("statusCode", http.StatusOK))

}
