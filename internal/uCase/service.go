package uCase

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ClientRepo interface {
	Create(ctx context.Context, info []byte) (string, error)
	Delete(ctx context.Context, uid string) error
}

type RequestIdempotencyKey interface {
	IsExist(ctx context.Context, uid string) (bool, error)
	Commit(ctx context.Context, uid string) error
}

type UseCase struct {
	logger      *zap.Logger
	tracer      trace.Tracer
	requestRepo RequestIdempotencyKey
	clientRepo  ClientRepo
}

func NewUseCase(logger *zap.Logger, clientRepo ClientRepo, requestRepo RequestIdempotencyKey) (*UseCase, error) {
	tracer := otel.Tracer("uCase")
	if tracer == nil {
		return nil, ErrTraceProviderIsNotSet
	}

	uCase := &UseCase{
		logger:      logger,
		tracer:      tracer,
		clientRepo:  clientRepo,
		requestRepo: requestRepo,
	}

	return uCase, nil
}
