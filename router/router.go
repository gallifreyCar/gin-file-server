package router

import (
	"encoding/json"
	"github.com/gallifreyCar/gin-file-server/handler"
	kafka_message "github.com/gallifreyCar/gin-file-server/kafka-message"
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func KafkaMiddleware(address []string, topic string) gin.HandlerFunc {

	return func(c *gin.Context) {

		//set a logger
		logFile, logger := log2.InitLogFile("gin-file-server.log", "[KafkaMiddleware]")
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logger.Println(err)
			}
		}(logFile)

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

			logger.Println(fieldName)
			// Parse the multipart form
			forms := c.Request.MultipartForm
			logger.Println(forms)
			// Get the uploaded files based on the specified field name
			files := forms.File[fieldName]
			for _, file := range files {
				logger.Println(file.Filename)
				logger.Println(file.Header)
				//Try to open the file.
				_, err := file.Open()
				if err != nil {
					// If an error occurred, return a bad request response.
					err = c.AbortWithError(http.StatusBadRequest, err)
					logger.Println(err)
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
		err := kafka_message.ProduceWriter(address, topic, kafkaMessages)
		if err != nil {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			logger.Println(err)
			return
		}
		// Continue processing the request.
		c.Next()
	}
}

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
	return r
}
