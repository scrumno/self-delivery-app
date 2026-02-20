package get_user_by_phone

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

type Fetcher struct {
	repository user.UserRepositoryInterface
}

func NewFetcher(repository user.UserRepositoryInterface) *Fetcher {
	return &Fetcher{
		repository: repository,
	}
}

func (f *Fetcher) Fetch(q Query) (UserDTO, error) {
	res, err := f.repository.GetByPhone(q.Phone)
	if err != nil {
		return UserDTO{}, err
	}

	return UserDTO{
		ID:        res.ID,
		Phone:     res.Phone,
		FullName:  res.FullName,
		BirthDate: res.BirthDate,
		IsActive:  res.IsActive,
		CreatedAt: res.CreatedAt,
	}, nil
}

type UserDTO struct {
	ID        uuid.UUID  `json:"id"`
	Phone     string     `json:"phone"`
	FullName  *string    `json:"full_name"`
	BirthDate *time.Time `json:"birth_date"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
}
