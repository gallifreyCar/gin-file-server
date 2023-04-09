package log

import (
	"log"
	"os"
)

func InitLogFile(fileName, prefix string) (*os.File, *log.Logger) {

	file, err := os.OpenFile("../log/"+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)

	}
	logger := log.New(file, prefix, log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetOutput(file)
	return file, logger
}
