package utils

import (
	"crypto/rand"
	"errors"
	mLogger "github.com/gallifreyCar/gin-file-server/m-logger"
	"go.uber.org/zap"
	"os"
)

func CreateTempFiles(fileNum int, fileSize int64, suffix string) (files []*os.File, close func(), err error) {
	// Validate args
	if fileNum <= 0 {
		return nil, nil, errors.New("fileNum must be greater than 0")
	}

	// Set a zap logger
	logger, err, closeFunc := mLogger.InitZapLogger("gin-file-server.log", "[CreateFiles]")
	logger.Error("Failed to init zap logger", zap.Error(err))
	defer closeFunc()

	// Create temporary files
	for i := 0; i < fileNum; i++ {
		file, err := os.CreateTemp("", suffix)
		if err != nil {
			logger.Error("Failed to create temporary files", zap.Error(err))
			return nil, nil, err
		}

		// Generate len(data) random bytes and writes them into data
		data := make([]byte, fileSize)
		_, err = rand.Read(data)
		if err != nil {
			logger.Error("Failed to generate len(data) random bytes and writes them into data", zap.Error(err))
			return nil, nil, err
		}
		// Write data into file
		if _, err := file.Write(data); err != nil {
			logger.Error("Failed to write data into file", zap.Error(err))
			return nil, nil, err
		}
		// Add file into files
		files = append(files, file)
	}

	// Remove files from os
	cleanFiles := func() {
		for _, file := range files {
			err := file.Close()
			if err != nil {
				logger.Error("Failed to create temporary files", zap.Error(err))
			}
			err = os.Remove(file.Name())
			if err != nil {
				logger.Error("Failed to create temporary files", zap.Error(err))
			}
		}
	}

	return files, cleanFiles, nil
}
