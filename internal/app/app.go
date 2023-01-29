package app

import (
	"context"
	"fmt"
	"github.com/Imm0bilize/gunshot-api-service/internal/config"
	"github.com/Imm0bilize/gunshot-api-service/internal/controller/grpc"
	"github.com/Imm0bilize/gunshot-api-service/internal/controller/http"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/msbroker"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	"github.com/Imm0bilize/gunshot-api-service/internal/uCase"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func createTraceProvider(cfg config.OTELConfig) func(context.Context) error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(net.JoinHostPort(cfg.Host, cfg.Port)),
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

func createKafkaProducer(cfg config.KafkaConfig) (sarama.SyncProducer, error) {
	kfkCfg := sarama.NewConfig()
	kfkCfg.Version = sarama.V2_5_0_0
	kfkCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(cfg.Peers, ","), kfkCfg)
	if err != nil {
		return nil, errors.Wrap(err, "error during create producer")
	}

	producer = otelsarama.WrapSyncProducer(kfkCfg, producer)

	return producer, nil
}

func createDB(cfg config.DBConfig) (*mongo.Database, func(context.Context) error, error) {
	clientOptions := options.Client()
	clientOptions.Monitor = otelmongo.NewMonitor()
	clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s:%s/", cfg.Host, cfg.Port))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error during create database")
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, errors.Wrap(err, "error during exec ping")
	}

	disconnect := func(ctx2 context.Context) error {
		return client.Disconnect(ctx)
	}

	return client.Database(cfg.Name), disconnect, nil
}

func Run(cfg *config.Config) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger.Debug("logger successfully created")

	ctx := context.Background()

	// Trace
	shutdownTraceProvider := createTraceProvider(cfg.OTEL)

	// DB
	db, dbShutdown, err := createDB(cfg.DB)
	logger.Debug("successfully connected to the database")

	producer, err := createKafkaProducer(cfg.Kafka)
	if err != nil {

	}
	// Broker
	broker := msbroker.NewKafkaProducer(logger, producer, cfg.Kafka.Topic)

	// domain service
	params := uCase.Params{
		Logger:      logger,
		Repo:        repository.NewRepo(db),
		AudioSender: broker,
		AudioLength: 1000,
	}

	useCase, err := uCase.NewUseCase(params)

	if err != nil {
		logger.Fatal("error when creating business logic of the service", zap.Error(err))
	}

	//http server
	httpServer := http.NewHTTPServer(logger, useCase)

	// grpc server
	listener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		logger.Fatal(err.Error())
	}

	//httpServer.Run(net.JoinHostPort("", cfg.HTTP.Port))
	go func() {
		if err := httpServer.Run(fmt.Sprintf(":%s", cfg.HTTP.Port)); err != nil {
			panic(err)
		}
	}()

	grpcServer := grpc.NewGRPCServer(logger, useCase)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("error during grpc server work", zap.Error(err))
		}
	}()

	// Shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	ctx, shutdownFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownFunc()

	if err = grpc.Shutdown(ctx, grpcServer); err != nil {
		logger.Error("error shutting down grpc server", zap.Error(err))
	}

	if err = dbShutdown(ctx); err != nil {
		logger.Error("error when closing database connection", zap.Error(err))
	}

	if err = broker.Shutdown(); err != nil {
		logger.Error("error when shutting down broker", zap.Error(err))
	}

	if err = shutdownTraceProvider(ctx); err != nil {
		logger.Error("error when shutting down provider", zap.Error(err))
	}
}
