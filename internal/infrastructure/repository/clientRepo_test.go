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
)

type ClientRepoSuite struct {
	suite.Suite
	db        *repository.ClientRepo
	client    *redis.Client
	container testcontainers.Container
}

func TestClientRepoSuite(t *testing.T) {
	suite.Run(t, new(ClientRepoSuite))
}

func (s *ClientRepoSuite) SetupSuite() {
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

	s.db = repository.NewClientRepo(s.client)
}

func (s *ClientRepoSuite) TestCreate() {
	uid := uuid.NewString()
	info := []byte(`{"info": "test"}`)

	err := s.db.Create(context.TODO(), uid, info)
	s.Nil(err)

	result, err := s.client.Get(context.TODO(), uid).Bytes()
	s.Equal(info, result)
}

func (s *ClientRepoSuite) TestDeletePositive() {
	uid := uuid.NewString()
	info := []byte(`{"info": "test"}`)

	err := s.db.Create(context.TODO(), uid, info)
	s.Nil(err)

	err = s.db.Delete(context.TODO(), uid)
	s.Nil(err)
}

func (s *ClientRepoSuite) TestDeleteNegative() {
	uid := uuid.NewString()
	err := s.db.Delete(context.TODO(), uid)
	s.Equal(err, repository.ErrClientNotFound)
}
