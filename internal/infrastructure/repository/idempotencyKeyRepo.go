package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
)

// RequestIdempotencyKeyRepo repository for idempotency keys
type RequestIdempotencyKeyRepo struct {
	db  *redis.Client
	ttl time.Duration
}

// IsExist checks for the presence of the key in the repository
func (r *RequestIdempotencyKeyRepo) IsExist(ctx context.Context, uid string) (bool, error) {
	n, err := r.db.Exists(ctx, uid).Result()
	if err != nil {
		return false, fmt.Errorf("error when getting data from the repository: %w", err)
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

// Commit writes the processed key to the database
func (r *RequestIdempotencyKeyRepo) Commit(ctx context.Context, uid string) error {
	return r.db.Set(ctx, uid, true, r.ttl).Err()
}

func NewIdempotencyKeyRepo(db *redis.Client, ttl time.Duration) *RequestIdempotencyKeyRepo {
	return &RequestIdempotencyKeyRepo{
		db:  db,
		ttl: ttl,
	}
}
