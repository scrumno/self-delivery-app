package factory

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository[T any] struct {
	DB *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{DB: db}
}

func (repo *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return repo.DB.WithContext(ctx).Create(entity).Error
}

func (repo *GormRepository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := repo.DB.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (repo *GormRepository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := repo.DB.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (repo *GormRepository[T]) FindAll(ctx context.Context, offset, limit int) ([]T, error) {
	var entities []T
	err := repo.DB.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
	return entities, err
}

func (repo *GormRepository[T]) GetAll(ctx context.Context, offset, limit int) ([]T, error) {
	var entities []T
	err := repo.DB.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
	return entities, err
}

func (repo *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return repo.DB.WithContext(ctx).Save(entity).Error
}

func (repo *GormRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	return repo.DB.WithContext(ctx).Delete(&entity, id).Error
}
