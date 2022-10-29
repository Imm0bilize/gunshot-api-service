package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ClientRepo struct {
	db     *redis.Client
	tracer trace.Tracer
}

// Create save a new user with uuid and information about him
func (c *ClientRepo) Create(ctx context.Context, uid string, info []byte) error {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Create")
	defer span.End()

	if err := c.db.Set(ctx, uid, info, 0).Err(); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to write the client to the database: %s", err.Error()))
		return fmt.Errorf("failed to write the client to the database: %w", err)
	}
	return nil
}

// Delete remove the user from database
func (c *ClientRepo) Delete(ctx context.Context, uid string) error {
	ctx, span := c.tracer.Start(ctx, "ClientRepo.Delete")
	defer span.End()

	n, err := c.db.Del(ctx, uid).Result()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to delete the user from the database: %w", err)
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
