package kafka_message

import (
	"context"
	"fmt"
	"github.com/gallifreyCar/gin-file-server/m-logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Consume(topic string) (message string, err error) {

	//set a zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server_"+time.Now().Format("20060102")+".log", "[Consume]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	// to consume messages
	partition := 0
	address := os.Getenv("address")

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		logger.Info("Failed to dial leader:", zap.Error(err))
		return "", err
	}

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		logger.Info("Failed to set ReadDeadline:", zap.Error(err))
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
			logger.Info("Failed to close batch:", zap.Error(err))
		}
		if err := conn.Close(); err != nil {
			logger.Info("Failed to close connection:", zap.Error(err))
		}
	}()
	return message, nil

}

func ConsumeReader(brokers []string, topic string, partition int, stop chan bool, done chan int) (err error) {

	//set a zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server_"+time.Now().Format("20060102")+".log", "[ConsumeReader]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	// make a new reader
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
		logger.Error("Failed to set kafka partition offset", zap.Error(err))
		return err
	}
	defer func(r *kafka.Reader) {
		err := r.Close()
		if err != nil {
			logger.Error("Failed to close kafka reader", zap.Error(err))

		}
	}(r)

	//create a file "consumerFile" and logger2 to log the topic have been consumed
	logger2, closeFunc := m_logger.InitLogFile("../target/consumer/"+time.Now().Format(time.DateOnly), "[🚄CONSUMING]")
	defer closeFunc()
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
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			msg, err := r.ReadMessage(ctx)
			cancel()

			if err != nil && err != context.DeadlineExceeded {
				logger2.Printf("Error reading message: %v\n", err)
				continue
			}
			logger2.Printf("Received message from topic %v and partition %d:  %s = %s\n", msg.Topic, msg.Partition, string(msg.Key), string(msg.Value))
		}
	}

}
