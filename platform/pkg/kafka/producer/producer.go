package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ platformKafka.Producer = (*Producer)(nil)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *logger.Logger
}

func New(brokers []string, topic string, log *logger.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("create Kafka producer: %w", err)
	}
	return &Producer{producer: syncProducer, topic: topic, logger: log}, nil
}

func (p *Producer) Send(ctx context.Context, key, value []byte) error {
	partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		return fmt.Errorf("send Kafka message: %w", err)
	}
	p.logger.Info(ctx, "Kafka message published",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)
	return nil
}

func (p *Producer) Close() error { return p.producer.Close() }
