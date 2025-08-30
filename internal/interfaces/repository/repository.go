package repository

import (
	"context"

	"github.com/google/uuid"
	"panel.go/internal/ent"
)

type BaseRepository struct {
	Ent *ent.Client
}

type BasePayload interface{}

type BaseResponse interface{}

type BaseRepositoryInterface[C BasePayload, U BasePayload, R BaseResponse] interface {
	FindOne(ctx context.Context, id uuid.UUID) (R, error)
	FindAll(ctx context.Context) ([]R, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, payload C) (R, error)
	Update(ctx context.Context, payload U) (R, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
