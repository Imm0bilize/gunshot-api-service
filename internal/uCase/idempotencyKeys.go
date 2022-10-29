package uCase

import (
	"context"
	"github.com/pkg/errors"
)

func (u *UseCase) checkRequestInIdempotencyRepo(ctx context.Context, key string) error {
	ctx, span := u.tracer.Start(ctx, "uCase.CheckRequestInIdempotencyRepo")
	defer span.End()

	if isExist, err := u.requestRepo.IsExist(ctx, key); err != nil {
		return errors.Wrap(err, "can't checks reqID in idempotency repository")
	} else {
		if isExist {
			return ErrRequestAlreadyProcessed
		}
	}

	return nil
}

func (u *UseCase) commitRequest(ctx context.Context, key string) error {
	ctx, span := u.tracer.Start(ctx, "uCase.CommitRequest")
	defer span.End()

	if err := u.requestRepo.Commit(ctx, key); err != nil {
		return errors.Wrap(err, "error when committing request in repository")
	}

	return nil
}
