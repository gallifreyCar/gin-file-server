package kafka_message

import (
	"context"
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/segmentio/kafka-go"
	"os"
	"time"
)

func Produce(topic string, message []string) (err error) {

	//set a logger
	logFile, log := log2.InitLogFile("gin-file-server.log", "[Produce]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			log.Fatal(err)
		}
	}(logFile)

	// to produce messages
	partition := 0
	address := os.Getenv("address")

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {

		log.Println("failed to dial leader:", err)
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {

		log.Fatal("failed to set WriteDeadline:", err)
		return err
	}

	for _, v := range message {
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(v)},
		)
	}

	if err != nil {

		log.Fatal("failed to write messages:", err)
		return err
	}

	if err := conn.Close(); err != nil {

		log.Fatal("failed to close writer:", err)
		return err
	}

	return nil
}
