package kafka_message

import (
	"context"
	"errors"
	"github.com/gallifreyCar/gin-file-server/m-logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"time"
)

// Produce use kafka-go Connection api, learn more:https://pkg.go.dev/github.com/segmentio/kafka-go#readme-connection
func Produce(topic string, message []string, partition int) (err error) {

	//set a zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server_"+time.Now().Format("20060102")+".log", "[Produce]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

	// to produce messages
	address := os.Getenv("address")

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		if err != nil {
			logger.Error("Failed to init zap logger", zap.Error(err))
		}
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {

		logger.Error("Failed to set WriteDeadline", zap.Error(err))
		return err
	}

	for _, v := range message {
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(v)},
		)
	}

	if err != nil {

		logger.Error("Failed to write messages", zap.Error(err))
		return err
	}

	if err := conn.Close(); err != nil {

		logger.Error("Failed to close writer", zap.Error(err))
		return err
	}

	return nil
}

// ProduceWriter use kafka-go writer api, learn more:https://pkg.go.dev/github.com/segmentio/kafka-go#readme-writer
func ProduceWriter(address []string, topic string, messages []kafka.Message) (err error) {
	//set a zap logger
	logger, err, closeFunc := m_logger.InitZapLogger("gin-file-server_"+time.Now().Format("20060102")+".log", "[ProduceWriter]")
	if err != nil {
		logger.Error("Failed to init zap logger", zap.Error(err))
	}
	defer closeFunc()

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
			logger.Error("Failed to create topic", zap.Error(err))
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			logger.Error("Unexpected error ", zap.Error(err))
			return err
		}
		break
	}

	if err := w.Close(); err != nil {
		logger.Error("Failed to close writer:", zap.Error(err))
		return err
	}
	return nil
}
