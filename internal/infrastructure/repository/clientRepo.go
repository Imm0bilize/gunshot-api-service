package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
)

type ClientRepo struct {
	db *redis.Client
}

// Create save a new user with uuid and information about him
func (c *ClientRepo) Create(ctx context.Context, uid string, info []byte) error {
	if err := c.db.Set(ctx, uid, info, 0).Err(); err != nil {
		return fmt.Errorf("failed to write the client to the database: %w", err)
	}
	return nil
}

// Delete remove the user from database
func (c *ClientRepo) Delete(ctx context.Context, uid string) error {
	n, err := c.db.Del(ctx, uid).Result()
	if err != nil {
		return fmt.Errorf("failed to delete the user from the database: %w", err)
	}

	if n == 0 {
		return ErrClientNotFound
	}

	return nil
}

func NewClientRepo(db *redis.Client) *ClientRepo {
	return &ClientRepo{
		db: db,
	}
}
