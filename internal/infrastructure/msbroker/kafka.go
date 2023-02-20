package msbroker

import (
	"context"
	"encoding/json"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/Imm0bilize/gunshot-api-service/pkg/api/brokerschemas"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type KafkaProducer struct {
	topic    string
	tracer   trace.Tracer
	producer sarama.SyncProducer
	logger   *zap.Logger
}

func NewKafkaProducer(logger *zap.Logger, producer sarama.SyncProducer, topic string) *KafkaProducer {
	tracer := otel.Tracer("msbroker")

	return &KafkaProducer{
		tracer:   tracer,
		producer: producer,
		logger:   logger,
		topic:    topic,
	}
}

func (k *KafkaProducer) Send(ctx context.Context, reqID uuid.UUID, message entities.AudioMessage) error {
	ctx, span := k.tracer.Start(ctx, "msbroker.Send")
	defer span.End()

	msg := brokerschemas.AudioMessage{
		RequestID: reqID,
		Payload:   message,
	}

	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		return errors.Wrap(err, "can't marshal msg")
	}

	producerMsg := &sarama.ProducerMessage{
		Topic:     k.topic,
		Key:       sarama.StringEncoder(reqID.String()),
		Value:     sarama.ByteEncoder(msgBytes),
		Timestamp: time.Now(),
	}

	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(producerMsg))

	partition, offset, err := k.producer.SendMessage(producerMsg)
	if err != nil {
		return errors.Wrap(err, "can't send message into kafka")
	}

	k.logger.Info(
		"message successfully send to broker",
		zap.String("requestID", reqID.String()),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)

	return nil
}

func (k *KafkaProducer) Shutdown() error {
	return k.producer.Close()
}
