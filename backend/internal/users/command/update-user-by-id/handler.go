package update_user_by_id

import (
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/users/command/update-user-by-id/dto"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

type Handler struct {
	repository user.UserRepositoryInterface
}

func NewHandler(repository user.UserRepositoryInterface) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Handle(cmd Command) (dto.UserDTO, error) {

	user, err := h.repository.UpdateById(cmd.ID, utils.BuildUpdateMap(cmd.Data))
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		ID:         user.ID,
		Phone:      user.Phone,
		FullName:   user.FullName,
		BirthDate:  user.BirthDate,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		StaffRoles: user.StaffRoles,
	}, nil
}
