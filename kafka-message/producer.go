package kafka_message

import (
	"context"
	"errors"
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/segmentio/kafka-go"

	"os"
	"time"
)

// Produce use kafka-go Connection api, learn more:https://pkg.go.dev/github.com/segmentio/kafka-go#readme-connection
func Produce(topic string, message []string, partition int) (err error) {

	//set a logger
	logFile, log := log2.InitLogFile("gin-file-server.log", "[Produce]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(logFile)

	// to produce messages
	address := os.Getenv("address")

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {

		log.Println("failed to dial leader:", err)
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {

		log.Println("failed to set WriteDeadline:", err)
		return err
	}

	for _, v := range message {
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(v)},
		)
	}

	if err != nil {

		log.Println("failed to write messages:", err)
		return err
	}

	if err := conn.Close(); err != nil {

		log.Println("failed to close writer:", err)
		return err
	}

	return nil
}

// ProduceWriter use kafka-go writer api, learn more:https://pkg.go.dev/github.com/segmentio/kafka-go#readme-writer
func ProduceWriter(address []string, topic string, messages []kafka.Message) (err error) {
	//set a logger
	logFile, log := log2.InitLogFile("gin-file-server.log", "[ProduceWriter]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(logFile)

	w := &kafka.Writer{
		Addr:                   kafka.TCP(address...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}

	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = w.WriteMessages(ctx, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			log.Printf("error creating topic: %v", err)
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Printf("unexpected error %v", err)
			return err
		}
		break
	}

	if err := w.Close(); err != nil {
		log.Println("failed to close writer:", err)
		return err
	}
	return nil
}
