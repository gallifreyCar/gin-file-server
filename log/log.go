package log

import (
	"log"
	"os"
)

func InitLogFile(fileName, prefix string) *os.File {

	file, err := os.OpenFile("../log/"+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)

	}

	log.SetPrefix(prefix)
	log.SetOutput(file)
	return file
}
