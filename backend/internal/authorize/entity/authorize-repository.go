package entity

import (
	"context"
	"errors"

	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type RegistrationRepository interface {
	base.BaseRepository[User]
	FindByPhone(ctx context.Context, phone string) (*User, error)
	UpdateFieldsByPhone(ctx context.Context, phone string, fields map[string]any) error
}

type registrationGormRepository struct {
	*factory.GormRepository[User]
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationGormRepository{
		GormRepository: factory.NewGormRepository[User](db),
	}
}

func (r *registrationGormRepository) Create(ctx context.Context, entity *User) error {
	return r.DB.WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(entity).Error
}

func (r *registrationGormRepository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var u User
	err := r.DB.WithContext(ctx).
		Where("phone = ?", phone).
		First(&u).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *registrationGormRepository) UpdateFieldsByPhone(ctx context.Context, phone string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	return r.DB.WithContext(ctx).
		Model(&User{}).
		Where("phone = ?", phone).
		Updates(fields).Error
}
