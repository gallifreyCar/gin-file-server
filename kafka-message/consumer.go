package kafka_message

import (
	"context"
	"fmt"
	log2 "github.com/gallifreyCar/gin-file-server/log"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Consume(topic string) (message string, err error) {

	//set a logger
	logFile, logger := log2.InitLogFile("gin-file-server.log", "[Consume]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

			logger.Fatal(err)
		}
	}(logFile)
	// to consume messages
	partition := 0
	address := os.Getenv("address")

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		logger.Println("failed to dial leader:", err)
		return "", err
	}

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		logger.Println("failed to set ReadDeadline:", err)
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
			logger.Println("failed to close batch:", err)
		}
		if err := conn.Close(); err != nil {
			logger.Println("failed to close connection:", err)
		}
	}()
	return message, nil

}

func ConsumeReader(brokers []string, topic string, partition int, stop chan bool, done chan int) (err error) {

	//set a logger
	logFile, logger := log2.InitLogFile("gin-file-server.log", "[ConsumeReader]")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}(logFile)

	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   2 * time.Second,
	})
	err = r.SetOffset(kafka.LastOffset)
	if err != nil {
		logger.Println(err)
		return err
	}
	defer func(r *kafka.Reader) {
		err := r.Close()
		if err != nil {
			logger.Println(err)

		}
	}(r)

	//create a file "consumerFile" and logger2 to log the topic have been consumed
	consumerFile, err := os.OpenFile("../target/consumer/"+time.Now().Format(time.DateOnly), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Println(err)
		return err
	}
	logger2 := log.New(consumerFile, "[🚄CONSUMER]", log.Ltime|log.Lshortfile)
	logger2.SetOutput(consumerFile)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Println(err)
		}
	}(consumerFile)

	//read message from Kafka topic and log it into "consumerFile"
	logger2.Println("Start consuming...")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-stop:
			done <- 200
			logger2.Println("Received stop signal. Stopping the consumer...")
			return

		case <-signals:
			logger2.Println("Interrupt signal received. Closing...")
			return

		default:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			msg, err := r.ReadMessage(ctx)
			cancel()
			if err != nil {
				logger2.Printf("Error reading message: %v\n", err)
				continue
			}
			logger2.Printf("Received message from partition %d:  %s = %s\n", msg.Partition, string(msg.Key), string(msg.Value))
		}
	}

}
