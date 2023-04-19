package utils

import (
	"bytes"
	mLogger "github.com/gallifreyCar/gin-file-server/m-logger"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"time"
)

func CreateUploadBody(field string, files []*os.File) (*bytes.Buffer, error) {

	// Set a zap logger
	logger, err, closeFunc := mLogger.InitZapLogger("util"+time.Now().Format(time.DateOnly)+".log", "[CreateUploadBody]")
	logger.Error("Failed to init zap logger", zap.Error(err))
	defer closeFunc()

	// Create a new request with the files as the body of the request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, file := range files {

		part, err := writer.CreateFormFile(field, file.Name())
		if err != nil {
			logger.Error("Failed to creates a form-data header", zap.Error(err))
			return nil, err
		}
		if _, err = file.Seek(0, 0); err != nil {
			logger.Error("Failed to seek file offer 0 form 0 ", zap.Error(err))
			return nil, err
		}
		if _, err = io.Copy(part, file); err != nil {
			logger.Error("Failed to copy file to part", zap.Error(err))
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		logger.Error("Failed to close writer", zap.Error(err))
		return nil, err
	}

	return body, err
}
