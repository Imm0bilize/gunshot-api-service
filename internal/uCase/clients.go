package uCase

import (
	"context"
	"encoding/json"
	"github.com/Imm0bilize/gunshot-api-service/internal/controller/http/dto"
	"github.com/pkg/errors"
)

func (u *UseCase) CreateNewClient(ctx context.Context, reqID string, information dto.ClientInfo) (string, error) {
	ctx, span := u.tracer.Start(ctx, "uCase.CreateNewClient")
	defer span.End()

	if err := u.checkRequestInIdempotencyRepo(ctx, reqID); err != nil {
		return "", err
	}

	bInfo, err := json.Marshal(information)
	if err != nil {
		return "", errors.Wrap(err, "error when marshal client information")
	}

	uid, err := u.clientRepo.Create(ctx, bInfo)
	if err != nil {
		return "", errors.Wrap(err, "error when create new client in repo")
	}

	if err := u.commitRequest(ctx, reqID); err != nil {
		return "", err
	}

	return uid, nil
}

func (u *UseCase) DeleteClient(ctx context.Context, reqID string, uid string) error {
	ctx, span := u.tracer.Start(ctx, "uCase.DeleteClient")
	defer span.End()

	if err := u.checkRequestInIdempotencyRepo(ctx, reqID); err != nil {
		return err
	}

	if err := u.clientRepo.Delete(ctx, uid); err != nil {
		return err
	}

	if err := u.commitRequest(ctx, reqID); err != nil {
		return err
	}

	return nil
}

func (u *UseCase) UploadAudio(ctx context.Context, reqID string) {
	ctx, span := u.tracer.Start(ctx, "uCase.UploadAudio")
	defer span.End()
}
