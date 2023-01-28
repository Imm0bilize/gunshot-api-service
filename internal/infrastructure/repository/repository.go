package repository

import (
	"context"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	_ ClientRepository = ClientRepo{}
)

type ClientRepository interface {
	Create(ctx context.Context, client *entities.Client) (string, error)
	Get(ctx context.Context, id string) (entities.Client, error)
	Update(ctx context.Context, id string, client *entities.Client) error
	Delete(ctx context.Context, id string) error
}

type Repo struct {
	Client ClientRepository
}

func NewRepo(database *mongo.Database) *Repo {
	return &Repo{
		Client: NewClientRepo(database),
	}
}
