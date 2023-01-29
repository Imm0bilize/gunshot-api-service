package uCase

import (
	"context"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ClientRepo interface {
	Create(ctx context.Context, client *entities.Client) (string, error)
	Get(ctx context.Context, id string) (entities.Client, error)
	Update(ctx context.Context, id string, client *entities.Client) error
	Delete(ctx context.Context, id string) error
}

type Client struct {
	tracer     trace.Tracer
	clientRepo ClientRepo
	logger     *zap.Logger
}

func NewClientUCase(logger *zap.Logger, clientRepo ClientRepo) *Client {
	return &Client{
		logger:     logger,
		tracer:     otel.Tracer("uCase.Client"),
		clientRepo: clientRepo,
	}
}

func (c Client) Create(ctx context.Context, reqID uuid.UUID, client *entities.Client) (string, error) {
	ctx, span := c.tracer.Start(ctx, "uCase.Client.Create")
	defer span.End()

	id, err := c.clientRepo.Create(ctx, client)
	if err != nil {
		c.logger.Error(
			"error during create new client",
			zap.String("reqID", reqID.String()),
			zap.Error(err),
		)

		return "", errors.Wrap(err, "can't create new client")
	}

	return id, nil
}

func (c Client) Get(ctx context.Context, reqID uuid.UUID, clientID string) (entities.Client, error) {
	ctx, span := c.tracer.Start(ctx, "uCase.Client.Get")
	defer span.End()

	client, err := c.clientRepo.Get(ctx, clientID)
	if err != nil {
		return entities.Client{}, errors.Wrap(err, "can't get the client")
	}

	return client, nil
}

func (c Client) Update(ctx context.Context, reqID uuid.UUID, clientID string, client *entities.Client) error {
	ctx, span := c.tracer.Start(ctx, "uCase.Client.Update")
	defer span.End()

	err := c.clientRepo.Update(ctx, clientID, client)
	if err != nil {
		return errors.Wrap(err, "can't update the client")
	}

	return nil
}

func (c Client) Delete(ctx context.Context, reqID uuid.UUID, clientID string) error {
	ctx, span := c.tracer.Start(ctx, "uCase.Client.Delete")
	defer span.End()

	if err := c.clientRepo.Delete(ctx, clientID); err != nil {
		return errors.Wrap(err, "can't delete the client")
	}

	return nil
}
