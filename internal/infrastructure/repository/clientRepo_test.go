//go:build integration
// +build integration

package repository_test

import (
	"context"
	"fmt"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const (
	_mongoImageName = "mongo:5.0.9"
	_dbName         = "GunshotService"
	_collectionName = "Clients"
)

type ClientRepoSuite struct {
	suite.Suite
	repo      *repository.ClientRepo
	dbClient  *mongo.Client
	container testcontainers.Container
}

var tempoClient = &entities.Client{
	ID:                  primitive.ObjectID{},
	LocationName:        "test",
	FullName:            "test test",
	Latitude:            52.124,
	Longitude:           12.235,
	NotificationMethods: nil,
}

func TestClientRepoSuite(t *testing.T) {
	suite.Run(t, new(ClientRepoSuite))
}

func (c *ClientRepoSuite) SetupSuite() {
	ctx := context.Background()

	port, err := nat.NewPort("", "27017")
	c.Require().NoError(err)

	req := testcontainers.ContainerRequest{
		Image:        _mongoImageName,
		ExposedPorts: []string{string(port)},
		WaitingFor:   wait.ForListeningPort(port),
	}

	mongoC, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	c.Require().NoError(err)

	c.container = mongoC
	endpoint, err := c.container.Endpoint(ctx, "")
	if err != nil {
		c.T().Fatal(err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", endpoint)))
	c.Require().NoError(err)

	c.dbClient = mongoClient
	c.repo = repository.NewClientRepo(mongoClient.Database(_dbName))
}

func (c *ClientRepoSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c.Require().NoError(c.dbClient.Disconnect(ctx))
	c.Require().NoError(c.container.Terminate(ctx))
}

func (c *ClientRepoSuite) TestCreate() {
	testTable := []struct {
		name   string
		client *entities.Client
	}{
		{
			name:   "default creating",
			client: tempoClient,
		},
	}

	for _, testCase := range testTable {
		c.Run(testCase.name, func() {
			id, err := c.repo.Create(context.Background(), testCase.client)
			c.Require().NoError(err)

			objectID, err := primitive.ObjectIDFromHex(id)
			c.Require().NoError(err)

			testCase.client.ID = objectID
			var gotClient *entities.Client

			err = c.dbClient.Database(_dbName).Collection(_collectionName).FindOne(
				context.Background(),
				bson.M{
					"_id": objectID,
				},
			).Decode(&gotClient)
			c.Require().NoError(err)

			c.Equal(testCase.client, gotClient)
		})
	}
}

func (c *ClientRepoSuite) TestGet() {
	testTable := []struct {
		name           string
		id             string
		createClientFn func() (string, error)
		expErr         error
	}{
		{
			name:   "incorrect id",
			id:     "test",
			expErr: primitive.ErrInvalidHex,
		},
		{
			name:   "no client with the id",
			id:     "507f191e810c19729de860ea",
			expErr: repository.ErrClientNotFound,
		},
		{
			name:   "existing id",
			expErr: nil,
			createClientFn: func() (string, error) {
				return c.repo.Create(context.Background(), tempoClient)
			},
		},
	}

	for _, testCase := range testTable {
		c.Run(testCase.name, func() {
			if testCase.createClientFn != nil {
				id, err := testCase.createClientFn()
				c.Require().NoError(err)
				testCase.id = id
			}
			_, err := c.repo.Get(context.Background(), testCase.id)
			c.ErrorIs(err, testCase.expErr)
		})
	}
}

//func (c *ClientRepoSuite) TestUpdate() {
//	testTable := []struct {
//		name string
//	}{
//		{},
//	}
//
//	for _, testCase := range testTable {
//		c.Run(testCase.name, func() {
//
//		})
//	}
//}

//
//func (c *ClientRepoSuite) TestDelete() {
//	testTable := []struct {
//		name string
//	}{
//		{},
//	}
//
//	for _, testCase := range testTable {
//		c.Run(testCase.name, func() {
//
//		})
//	}
//}
