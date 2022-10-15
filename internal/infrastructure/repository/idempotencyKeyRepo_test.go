package repository_test

import (
	"context"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

const redisImageName = "redis:7.0-alpine3.16"

type IdempotenceKeySuite struct {
	suite.Suite
	db        *repository.RequestIdempotencyKeyRepo
	client    *redis.Client
	container testcontainers.Container
}

func TestIdempotenceKeySuite(t *testing.T) {
	suite.Run(t, new(IdempotenceKeySuite))
}

func (s *IdempotenceKeySuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        redisImageName,
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		s.T().Fatal(err)
	}

	s.container = redisC

	endpoint, err := s.container.Endpoint(ctx, "")
	if err != nil {
		s.T().Fatal(err)
	}

	s.client = redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	s.db = repository.NewIdempotencyKeyRepo(s.client, time.Second*20)
}

func (s *IdempotenceKeySuite) TestIsExist() {
	testTable := []struct {
		name         string
		uid          string
		insertDataFn func(ctx context.Context, uid string) error
		expectedRes  bool
		expectedErr  error
	}{
		{
			name: "default usage",
			uid:  uuid.NewString(),
			insertDataFn: func(ctx context.Context, uid string) error {
				return s.client.Set(ctx, uid, true, time.Second*10).Err()
			},
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "no data in db",
			uid:  uuid.NewString(),
			insertDataFn: func(ctx context.Context, uid string) error {
				return nil
			},
			expectedRes: false,
			expectedErr: nil,
		},
	}

	for _, testCase := range testTable {
		s.Run(testCase.name, func() {
			ctx := context.TODO()

			if err := testCase.insertDataFn(ctx, testCase.uid); err != nil {
				s.T().Error(err)
			}

			got, err := s.db.IsExist(ctx, testCase.uid)
			s.Equal(testCase.expectedErr, err)
			s.Equal(testCase.expectedRes, got)
		})
	}
}

func (s *IdempotenceKeySuite) TestIsExistDbDoesntWork() {
	dur := time.Second
	if err := s.container.Stop(context.TODO(), &dur); err != nil {
		s.T().Fatalf("failed to stop the container: %s", err.Error())
	}

	defer func() {
		if err := s.container.Start(context.TODO()); err != nil {
			s.T().Fatalf("failed to restart the container: %s", err.Error())
		}
	}()

	_, err := s.db.IsExist(context.TODO(), "test")
	s.NotNil(err)
}

func (s *IdempotenceKeySuite) TestCommit() {
	testTable := []struct {
		name        string
		uid         string
		getFromDB   func(ctx context.Context, uid string) (bool, error)
		expectedRes bool
		expectedErr error
	}{
		{
			name: "default usage",
			uid:  uuid.NewString(),
			getFromDB: func(ctx context.Context, uid string) (bool, error) {
				return s.client.Get(ctx, uid).Bool()
			},
			expectedRes: true,
			expectedErr: nil,
		},
	}

	for _, testCase := range testTable {
		s.Run(testCase.name, func() {
			ctx := context.TODO()

			err := s.db.Commit(ctx, testCase.uid)

			s.Equal(testCase.expectedErr, err)

			got, err := testCase.getFromDB(ctx, testCase.uid)
			s.Nil(err)

			s.Equal(testCase.expectedRes, got)
		})
	}
}

func (s *IdempotenceKeySuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.client.Close(); err != nil {
		s.T().Fatal(err)
	}

	if err := s.container.Terminate(ctx); err != nil {
		s.T().Fatal(err)
	}
}
