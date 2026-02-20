package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetByID(ID string) (User, error)
	GetByPhone(phone string) (User, error)
	GetAllUsers() ([]User, error)
	Create(user User) (User, error)
	Delete(ID string) error
	UpdateById(ID uuid.UUID, data map[string]interface{}) (User, error)
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

func (r *UserRepository) GetByPhone(phone string) (User, error) {
	var u User
	result := r.db.First(&u, "phone = ?", phone)
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}

func (r *UserRepository) GetAllUsers() ([]User, error) {

	var users []User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *UserRepository) Create(u User) (User, error) {
	result := r.db.Create(&u)
	if result.Error != nil {
		return User{}, result.Error
	}
	return u, nil
}

func (r *UserRepository) Delete(ID string) error {
	result := r.db.Model(&User{}).Where("id = ?", ID).Update("is_active", false)
	return result.Error
}

func (r *UserRepository) UpdateById(ID uuid.UUID, data map[string]interface{}) (User, error) {
	var u User

	if err := r.db.First(&u, "id = ?", ID).Error; err != nil {
		return User{}, err
	}

	if err := r.db.Model(&u).Updates(data).Error; err != nil {
		return User{}, err
	}

	return u, nil
}
