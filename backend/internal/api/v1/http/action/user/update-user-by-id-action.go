package user

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	updateUserById "github.com/scrumno/scrumno-api/internal/users/command/update-user-by-id"
)

type UpdateUserByIdAction struct {
	handler *updateUserById.Handler
}

func NewUpdateUserByIdAction(handler *updateUserById.Handler) *UpdateUserByIdAction {
	return &UpdateUserByIdAction{handler: handler}
}

type RequestData struct {
	Phone       string     `json:"phone"`
	FullName    *string    `json:"full_name,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	IikoGuestID *uuid.UUID `json:"iiko_guest_id,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (a *UpdateUserByIdAction) Action(w http.ResponseWriter, r *http.Request) {

	req, err := utils.DecodeUpdateRequest[RequestData](r)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	cmd := updateUserById.Command{
		ID:   req.ID,
		Data: updateUserById.CommandData{},
	}

	if req.IsSet("phone") {
		cmd.Data.Phone = &req.Fields.Phone
	}
	if req.IsSet("full_name") {
		cmd.Data.FullName = req.Fields.FullName
	}
	if req.IsSet("birth_date") {
		cmd.Data.BirthDate = req.Fields.BirthDate
	}
	if req.IsSet("iiko_guest_id") {
		cmd.Data.IikoGuestID = req.Fields.IikoGuestID
	}
	if req.IsSet("is_active") {
		cmd.Data.IsActive = &req.Fields.IsActive
	}
	if req.IsSet("created_at") {
		cmd.Data.CreatedAt = &req.Fields.CreatedAt
	}

	res, err := a.handler.Handle(cmd)
	if err != nil {
		utils.JSONResponse(w, map[string]string{"error": "Извините что-то пошло не так"}, 500)
		return
	}

	utils.JSONResponse(w, res, http.StatusOK)
}
