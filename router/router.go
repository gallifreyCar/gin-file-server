package router

import (
	"encoding/json"
	"github.com/gallifreyCar/gin-file-server/handler"
	kafkaMessage "github.com/gallifreyCar/gin-file-server/kafka-message"
	"github.com/gallifreyCar/gin-file-server/m-logger"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func KafkaMiddleware(address []string, topic string) gin.HandlerFunc {

	return func(c *gin.Context) {

		//set a zap logger
		logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server_"+time.Now().Format("20060102")+".log", "[KafkaMiddleware]")
		if err != nil {
			logger.Error("Failed to init zap logger", zap.Error(err))
		}
		defer closeFunc()

		method := c.Request.Method
		kafkaMessages := make([]kafka.Message, 0)
		if method == http.MethodGet {
			folder := c.Param("folder")
			fileName := c.Param("file_name")

			// Update file download metadata
			fileInfo := map[string]interface{}{
				"FilePath":       path.Join("..", "target", "upload", folder),
				"FileName":       fileName,
				"DownloadedTime": time.Now().Format(time.DateTime),
			}
			// Convert the file info to JSON.
			message, err := json.Marshal(fileInfo)
			if err != nil {
				log.Print(err)
			}
			// Construct the Kafka messages.
			kafkaMessages = []kafka.Message{{
				Key:   []byte(fileName),
				Value: message,
			}}

		}
		if method == http.MethodPost {
			//handle multiple files
			fieldName := c.DefaultPostForm("fieldName", "files")
			fileCount := len(c.Request.MultipartForm.File)

			//handle single file
			if fileCount == 1 {
				fieldName = c.DefaultPostForm("fieldName", "file")
			}

			// Parse the multipart form
			forms := c.Request.MultipartForm

			// Get the uploaded files based on the specified field name
			files := forms.File[fieldName]
			for _, file := range files {

				//Try to open the file.
				_, err := file.Open()
				if err != nil {
					// If an error occurred, return a bad request response.
					err = c.AbortWithError(http.StatusBadRequest, err)
					logger.Error("Fail to open file", zap.Error(err), zap.Int("statusCode", http.StatusBadRequest))
					return
				}

				// Create the file info map.
				fileInfo := map[string]interface{}{
					"FileName":   file.Filename,
					"MIMEHeader": file.Header,
					"FileSize":   file.Size,
					"UploadTime": time.Now().Format(time.DateTime),
				}
				// Convert the file info to JSON.
				message, err := json.Marshal(fileInfo)
				if err != nil {
					log.Print(err)
				}
				// Construct the Kafka messages.
				kafkaMessages = append(kafkaMessages, kafka.Message{
					Key:   []byte(file.Filename),
					Value: message,
				})

			}

		}

		// Write the message to Kafka.
		err = kafkaMessage.ProduceWriter(address, topic, kafkaMessages)
		if err != nil {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			logger.Error("Fail to write the message to Kafka", zap.Error(err), zap.Int("statusCode", http.StatusInternalServerError))
			return
		}

		logger.Info("Write the message to Kafka success!", zap.String("topic", topic))
		// Continue processing the request.
		c.Next()
	}
}

// MaxAllowed Middleware function to set a size limit on the uploaded files
func MaxAllowed(n int64) gin.HandlerFunc {
	//set a zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server_"+time.Now().Format("20060102")+".log", "[MaxAllowed]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	maxBytes := n
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		if err := c.Request.ParseMultipartForm(maxBytes); err != nil {
			if err.Error() == "http: request body too large" {
				c.AbortWithStatus(http.StatusRequestEntityTooLarge)
				logger.Error("Upload File is too large", zap.Error(err), zap.Int("statusCode", http.StatusRequestEntityTooLarge), zap.Int64("maxBytes", maxBytes))
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			logger.Error("Fail to parse multipart form in max bytes limit", zap.Error(err), zap.Int("statusCode", http.StatusBadRequest), zap.Int64("maxBytes", maxBytes))
			return
		}
		c.Next()
		logger.Info("File size check success!", zap.Int64("maxBytes", maxBytes))

	}

}

func SetupRouter() *gin.Engine {

	address1 := os.Getenv("address1")
	address2 := os.Getenv("address2")
	address := []string{address1, address2}

	r := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	//set a size limit on the uploaded files
	r.POST("/uploadFile", MaxAllowed(10<<20), handler.UploadFileSingle, KafkaMiddleware(address, "upload-file-single"))
	r.POST("/uploadFiles", MaxAllowed(50<<20), handler.UploadFiles, KafkaMiddleware(address, "upload-file-multiple"))
	r.GET("/downloadFile/:folder/:file_name", handler.DownloadFile, KafkaMiddleware(address, "download-file"))
	r.GET("/uploads/file_name/:file_name", handler.SelectFileLogByName)
	r.Use()
	return r
}
