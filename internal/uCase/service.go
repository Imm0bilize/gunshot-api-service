package uCase

import (
	"context"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	_ ClientUseCase = Client{}
	_ AudioUseCase  = Audio{}
)

type ClientUseCase interface {
	Create(ctx context.Context, reqID uuid.UUID, client *entities.Client) (string, error)
	Get(ctx context.Context, reqID uuid.UUID, id string) (entities.Client, error)
	Update(ctx context.Context, reqID uuid.UUID, id string, client *entities.Client) error
	Delete(ctx context.Context, reqID uuid.UUID, id string) error
}

type AudioUseCase interface {
	Upload(ctx context.Context, reqID uuid.UUID, id string, msg entities.AudioMessage) error
}

type UseCase struct {
	Client ClientUseCase
	Audio  AudioUseCase
}

type Params struct {
	Logger      *zap.Logger
	Repo        *repository.Repo
	AudioSender Sender
	AudioLength int
}

func NewUseCase(params Params) (*UseCase, error) {
	return &UseCase{
		Client: NewClientUCase(params.Logger, params.Repo.Client),
		Audio:  NewAudioUCase(params.Logger, params.AudioSender, params.AudioLength),
	}, nil
}
