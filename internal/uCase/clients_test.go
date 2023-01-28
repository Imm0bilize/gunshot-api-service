package uCase_test

import (
	"context"
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository"
	mock_repository "github.com/Imm0bilize/gunshot-api-service/internal/infrastructure/repository/mocks"
	"github.com/Imm0bilize/gunshot-api-service/internal/uCase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"testing"
)

func TestClientCreate(t *testing.T) {
	testTable := []struct {
		name          string
		setMockOutput func(context.Context, *entities.Client, *mock_repository.MockClientRepository)
		expErr        error
	}{
		{
			name:   "successfully creating",
			expErr: nil,
			setMockOutput: func(ctx context.Context, client *entities.Client, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Create")
				repo.EXPECT().Create(ctx, client).Return("", nil).Times(1)
			},
		},
		{
			name:   "db client disconnect",
			expErr: errors.New("can't create new client: client is disconnected"),
			setMockOutput: func(ctx context.Context, client *entities.Client, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Create")
				repo.EXPECT().Create(ctx, client).Return("", mongo.ErrClientDisconnected).Times(1)
			},
		},
	}

	for _, tCase := range testTable {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				ctx  = context.Background()
				ctrl = gomock.NewController(t)
				repo = mock_repository.NewMockClientRepository(ctrl)
			)

			client := &entities.Client{}
			tCase.setMockOutput(ctx, client, repo)

			useCase := uCase.NewClientUCase(zap.NewExample(), repo)

			_, err := useCase.Create(ctx, uuid.New(), client)

			if err != nil {
				require.Equal(t, tCase.expErr.Error(), err.Error())
			} else {
				require.ErrorIs(t, tCase.expErr, err)
			}

			ctrl.Finish()
		})
	}
}

func TestClientGet(t *testing.T) {
	testTable := []struct {
		name          string
		id            string
		expErr        error
		setMockOutput func(context.Context, string, *mock_repository.MockClientRepository)
	}{
		{
			name:   "successfully getting",
			id:     "123",
			expErr: nil,
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Get")
				repo.EXPECT().Get(ctx, id).Return(entities.Client{}, nil).Times(1)
			},
		},
		{
			name:   "db client disconnect",
			id:     "123",
			expErr: errors.New("can't get the client: client is disconnected"),
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Get")
				repo.EXPECT().Get(ctx, id).Return(entities.Client{}, mongo.ErrClientDisconnected).Times(1)
			},
		},
		{
			name:   "client not found",
			id:     "12",
			expErr: errors.New("can't get the client: the client is not found"),
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Get")
				repo.EXPECT().Get(ctx, id).Return(entities.Client{}, repository.ErrClientNotFound).Times(1)
			},
		},
	}
	for _, tCase := range testTable {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				ctx  = context.Background()
				ctrl = gomock.NewController(t)
				repo = mock_repository.NewMockClientRepository(ctrl)
			)

			tCase.setMockOutput(ctx, tCase.id, repo)

			useCase := uCase.NewClientUCase(zap.NewExample(), repo)
			_, err := useCase.Get(ctx, uuid.New(), tCase.id)

			if tCase.expErr != nil {
				require.Equal(t, err.Error(), tCase.expErr.Error())
			} else {
				require.ErrorIs(t, err, tCase.expErr)
			}

			ctrl.Finish()
		})
	}
}

func TestClientUpdate(t *testing.T) {
	testTable := []struct {
		name          string
		id            string
		expErr        error
		setMockOutput func(context.Context, string, *mock_repository.MockClientRepository)
	}{
		{
			name:   "successfully updating",
			id:     "123",
			expErr: nil,
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Update")
				repo.EXPECT().Update(ctx, id, &entities.Client{}).Return(nil).Times(1)
			},
		},
		{
			name:   "db client disconnect",
			id:     "123",
			expErr: errors.New("can't update the client: client is disconnected"),
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Update")
				repo.EXPECT().Update(ctx, id, &entities.Client{}).Return(mongo.ErrClientDisconnected).Times(1)
			},
		},
		{
			name:   "client not found",
			id:     "12",
			expErr: errors.New("can't update the client: the client is not found"),
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Update")
				repo.EXPECT().Update(ctx, id, &entities.Client{}).Return(repository.ErrClientNotFound).Times(1)
			},
		},
	}
	for _, tCase := range testTable {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				ctx  = context.Background()
				ctrl = gomock.NewController(t)
				repo = mock_repository.NewMockClientRepository(ctrl)
			)

			tCase.setMockOutput(ctx, tCase.id, repo)

			useCase := uCase.NewClientUCase(zap.NewExample(), repo)
			err := useCase.Update(ctx, uuid.New(), tCase.id, &entities.Client{})

			if tCase.expErr != nil {
				require.Equal(t, err.Error(), tCase.expErr.Error())
			} else {
				require.ErrorIs(t, err, tCase.expErr)
			}

			ctrl.Finish()
		})
	}
}

func TestClientDelete(t *testing.T) {
	testTable := []struct {
		name          string
		id            string
		expErr        error
		setMockOutput func(context.Context, string, *mock_repository.MockClientRepository)
	}{
		{
			name:   "successfully updating",
			id:     "123",
			expErr: nil,
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Delete")
				repo.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
		},
		{
			name:   "db client disconnect",
			id:     "123",
			expErr: errors.New("can't delete the client: client is disconnected"),
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Delete")
				repo.EXPECT().Delete(ctx, id).Return(mongo.ErrClientDisconnected).Times(1)
			},
		},
		{
			name:   "client not found",
			id:     "12",
			expErr: errors.New("can't delete the client: the client is not found"),
			setMockOutput: func(ctx context.Context, id string, repo *mock_repository.MockClientRepository) {
				ctx, _ = otel.GetTracerProvider().Tracer("uCase.Client").Start(ctx, "uCase.Client.Delete")
				repo.EXPECT().Delete(ctx, id).Return(repository.ErrClientNotFound).Times(1)
			},
		},
	}
	for _, tCase := range testTable {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				ctx  = context.Background()
				ctrl = gomock.NewController(t)
				repo = mock_repository.NewMockClientRepository(ctrl)
			)

			tCase.setMockOutput(ctx, tCase.id, repo)

			useCase := uCase.NewClientUCase(zap.NewExample(), repo)
			err := useCase.Delete(ctx, uuid.New(), tCase.id)

			if tCase.expErr != nil {
				require.Equal(t, err.Error(), tCase.expErr.Error())
			} else {
				require.ErrorIs(t, err, tCase.expErr)
			}

			ctrl.Finish()
		})
	}
}
