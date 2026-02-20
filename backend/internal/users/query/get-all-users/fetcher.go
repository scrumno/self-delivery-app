package get_all_users

import (
	"time"

	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

type Fetcher struct {
	repository user.UserRepositoryInterface
}

func NewFetcher(repository user.UserRepositoryInterface) *Fetcher {
	return &Fetcher{repository: repository}
}

func (f *Fetcher) Fetch(_ Query) (UsersDTO, error) {
	users, err := f.repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	dto := make(UsersDTO, 0, len(users))
	for _, u := range users {
		dto = append(dto, UserDTO{
			ID:        u.ID,
			Phone:     u.Phone,
			FullName:  u.FullName,
			BirthDate: u.BirthDate,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		})
	}

	return dto, nil
}

type UserDTO struct {
	ID        interface{} `json:"id"`
	Phone     string      `json:"phone"`
	FullName  *string     `json:"full_name,omitempty"`
	BirthDate interface{} `json:"birth_date,omitempty"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
}

type UsersDTO []UserDTO
