package delete_user

import (
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

type Handler struct {
	repository user.UserRepositoryInterface
}

func NewHandler(repository user.UserRepositoryInterface) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Handle(cmd Command) error {
	return h.repository.Delete(cmd.ID)
}
