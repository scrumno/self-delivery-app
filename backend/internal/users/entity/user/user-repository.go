package user

import "gorm.io/gorm"

type UserRepositoryInterface interface {
	GetByID(ID string) (User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func (r *UserRepository) GetByID(ID string) (User, error) {
	var user User

	result := r.db.First(&user, "id = ?", ID)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}
