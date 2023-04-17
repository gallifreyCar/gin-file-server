package kafka_message

import (
	"context"
	"fmt"
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/segmentio/kafka-go"
	"os"
	"time"
)

func Consume(topic string) (message string, err error) {

	//set a logger
	logFile, log := log2.InitLogFile("gin-file-server.log", "[Consume]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(logFile)
	// to consume messages
	partition := 0
	address := os.Getenv("address")

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		log.Println("failed to dial leader:", err)
		return "", err
	}

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Println("failed to set ReadDeadline:", err)
		return "", err
	}
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	message = "message: "
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
		message = message + string(b[:n])
	}

	defer func() {
		if err := batch.Close(); err != nil {
			log.Println("failed to close batch:", err)
		}
		if err := conn.Close(); err != nil {
			log.Println("failed to close connection:", err)
		}
	}()
	return message, nil

}
