package base

import (
	"context"

	"github.com/google/uuid"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, offset, limit int) ([]T, error)
	FindAll(ctx context.Context, offset, limit int) ([]T, error)
}
