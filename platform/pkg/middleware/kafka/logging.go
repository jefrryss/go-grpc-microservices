package kafka

import (
	"context"

	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/consumer"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	"go.uber.org/zap"
)

func Logging(log *logger.Logger) consumer.Middleware {
	return func(next platformKafka.Handler) platformKafka.Handler {
		return func(ctx context.Context, message platformKafka.Message) error {
			log.Info(ctx, "Kafka message received", zap.String("topic", message.Topic), zap.Int64("offset", message.Offset))
			return next(ctx, message)
		}
	}
}
