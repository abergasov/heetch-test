package broker

import (
	"context"
	"zombie_locator/internal/logger"

	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer defines a Kafka messages consumer.
type KafkaConsumer struct {
	exitMark chan struct{}
	log      logger.AppLogger
	reader   *kafka.Reader
	dlq      Producer
	topic    string
}

// NewKafkaConsumer sets up a new kafka consumer for the given topic using the provided handler.
func NewKafkaConsumer(log logger.AppLogger, dlq Producer, brokerAddrs []string, consumerGroupID, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		exitMark: make(chan struct{}),
		dlq:      dlq,
		log: log.With(zap.String("service", "kafka_consumer")).
			With(zap.String("group_id", consumerGroupID)).
			With(zap.String("topic", topic)),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokerAddrs,
			GroupID: consumerGroupID,
			Topic:   topic,
		}),
		topic: topic,
	}
}

// Run starts the kafka consumer.
func (p *KafkaConsumer) Run(ctx context.Context, handler Handler) (err error) {
	defer close(p.exitMark)
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			// Get the next message to consume from the broker
			m, err := p.reader.FetchMessage(ctx)
			if err != nil {
				return err
			}

			err = handler(ctx, m.Value)
			if err == nil {
				continue
			}
			if err = p.putInDeadLetter(ctx, err, m.Value); err != nil {
				p.log.Error("failed to handle message", err)
			}

			// Mark the processed message as committed on the broker.
			if err = p.reader.CommitMessages(ctx, m); err != nil {
				break loop
			}
		}
	}
	return err
}

// DeadLetter put failed message to dead letter queue. extra logic can be added here - retry, etc.
// message can store into db|cache
func (p *KafkaConsumer) putInDeadLetter(ctx context.Context, err error, msg []byte) error {
	return p.dlq.WriteDeadMessages(ctx, p.topic, err, msg)
}

func (p *KafkaConsumer) Shutdown() error {
	<-p.exitMark
	if err := p.dlq.Shutdown(); err != nil {
		p.log.Error("failed to shutdown dead letter queue", err)
	}
	return p.reader.Close()
}
