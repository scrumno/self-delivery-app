package status

import (
	"gorm.io/gorm"
)

type StatusRepositoryInterface interface {
	CheckStatus() error
}

type StatusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

var _ StatusRepositoryInterface = (*StatusRepository)(nil)

func (r *StatusRepository) CheckStatus() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
