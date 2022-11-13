package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zombie_locator/internal/logger"

	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

type deadMessage struct {
	Topic string `json:"topic"`
	Error string `json:"error"`
	Msg   []byte `json:"msg"`
}

// KafkaProducer defines a Kafka messages producer.
type KafkaProducer struct {
	log    logger.AppLogger
	writer *kafka.Writer
}

func NewKafkaProducer(log logger.AppLogger, brokerAddr string, topic string) *KafkaProducer {
	return &KafkaProducer{
		log: log.With(zap.String("component", "kafka_producer")).
			With(zap.String("topic", topic)),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokerAddr),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (k *KafkaProducer) WriteDeadMessages(ctx context.Context, srcTopic string, err error, msg []byte) error {
	d := deadMessage{
		Error: err.Error(),
		Topic: srcTopic,
		Msg:   msg,
	}
	data, err := json.Marshal(&d)
	if err != nil {
		return fmt.Errorf("failed to marshal dead message: %w", err)
	}
	for i := 0; i < 3; i++ {
		if err = k.writer.WriteMessages(ctx, kafka.Message{Value: data}); err == nil {
			return nil
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	k.log.Error("failed to write dead message", err)
	return err
}

func (k *KafkaProducer) Shutdown() error {
	// simply close connection, as dlq writes from consumer
	return k.writer.Close()
}
