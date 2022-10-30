package repository

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ClientRepo struct {
	db     *redis.Client
	tracer trace.Tracer
}

// Create save a new user with uuid and information about him
func (c *ClientRepo) Create(ctx context.Context, info []byte) (string, error) {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Create")
	defer span.End()

	uid := uuid.NewString()

	if err := c.db.Set(ctx, uid, info, 0).Err(); err != nil {
		rErr := errors.Wrap(err, "failed to write the client to the database")
		span.SetStatus(codes.Error, rErr.Error())
		return "", rErr
	}
	return uid, nil
}

func (c *ClientRepo) Update(ctx context.Context, uid string, info []byte) error {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Update")
	defer span.End()

	n, err := c.db.Exists(ctx, uid).Result()
	if err != nil {
		rErr := errors.Wrap(err, "failed to check existing the client in db")
		span.SetStatus(codes.Error, rErr.Error())
		return rErr
	}

	if n == 0 {
		return ErrClientNotFound
	}

	if err := c.db.Set(ctx, uid, info, 0).Err(); err != nil {
		rErr := errors.Wrap(err, "failed to update the client information")
		span.SetStatus(codes.Error, rErr.Error())
		return rErr
	}

	return nil
}

// Delete remove the user from database
func (c *ClientRepo) Delete(ctx context.Context, uid string) error {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Delete")
	defer span.End()

	n, err := c.db.Del(ctx, uid).Result()
	if err != nil {
		rErr := errors.Wrap(err, "failed to delete the user from the database")
		span.SetStatus(codes.Error, rErr.Error())
		return rErr
	}

	if n == 0 {
		return ErrClientNotFound
	}

	return nil
}

func NewClientRepo(db *redis.Client) *ClientRepo {
	tracer := otel.Tracer("ClientRepo")
	return &ClientRepo{
		db:     db,
		tracer: tracer,
	}
}
