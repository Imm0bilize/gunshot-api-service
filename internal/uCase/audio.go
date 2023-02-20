package uCase

import (
	"context"
	"fmt"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Sender interface {
	Send(ctx context.Context, reqID uuid.UUID, msg entities.AudioMessage) error
}

type Audio struct {
	audioSender Sender
	tracer      trace.Tracer
	logger      *zap.Logger
	audioLength int
}

var (
	ErrNotEqRequiredLength = errors.New("the audio not equal to the required length")
)

func NewAudioUCase(logger *zap.Logger, audioSender Sender, audioLength int) *Audio {
	return &Audio{
		audioSender: audioSender,
		tracer:      otel.Tracer("uCase.Audio"),
		audioLength: audioLength,
		logger:      logger,
	}
}

func (a Audio) Upload(ctx context.Context, reqID uuid.UUID, clientID string, msg entities.AudioMessage) error {
	ctx, span := a.tracer.Start(ctx, "uCase.Audio.Upload")
	defer span.End()

	castedID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return errors.Wrap(err, "error during convert client id")
	}
	msg.ID = castedID

	if err := a.validate(msg.Payload); err != nil {
		return errors.Wrap(err, "validation error")
	}

	if err := a.audioSender.Send(ctx, reqID, msg); err != nil {
		span.RecordError(err)
		return errors.Wrap(err, "error during send audio")
	}

	return nil
}

func (a Audio) validate(audio []byte) error {
	if len(audio) != a.audioLength {
		return fmt.Errorf(
			"%w: expected %d (bytes), got %d (bytes)", ErrNotEqRequiredLength, a.audioLength, len(audio),
		)
	}

	return nil
}
