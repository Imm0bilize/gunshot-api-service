package app

import (
	"context"
	"github.com/Imm0bilize/gunshot-api-service/internal/config"
	"github.com/Imm0bilize/gunshot-api-service/internal/controller/http"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	"github.com/Imm0bilize/gunshot-api-service/internal/uCase"
	"github.com/go-redis/redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"log"
	"net"
)

func createTraceProvider(collectorURL string) func(context.Context) error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "gunshot-api-service"),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}

func Run(cfg *config.Config) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger.Debug("logger successfully created")

	ctx := context.Background()

	// Trace
	shutdownTraceProvider := createTraceProvider(
		net.JoinHostPort(cfg.OTEL.Host, cfg.OTEL.Port),
	)

	defer func(ctx context.Context) {
		if err := shutdownTraceProvider(ctx); err != nil {
			logger.Error(err.Error())
		}
	}(ctx)

	// DB
	db := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.DB.Host, cfg.DB.Port),
		Password: cfg.DB.Password,
		DB:       0,
	})
	if err := db.Ping(ctx).Err(); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("successfully connected to the database")

	// domain service
	useCase, err := uCase.NewUseCase(
		logger,
		repository.NewClientRepo(db),
		repository.NewIdempotencyKeyRepo(db, cfg.IdemKey.TTL),
	)

	if err != nil {
		logger.Fatal(err.Error())
	}

	// http server
	httpServer, err := http.New(logger, useCase)
	if err != nil {
		logger.Fatal(err.Error())
	}

	httpServer.Run(net.JoinHostPort("", cfg.HTTP.Port))
	// grpc server

	// shutdown
}
