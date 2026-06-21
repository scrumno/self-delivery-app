package user

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/middleware"
	updateUserProfile "github.com/scrumno/scrumno-api/internal/users/command/update-user-profile"
)

type UpdateUserProfileAction struct {
	Handler *updateUserProfile.Handler
}

func NewUpdateUserProfileAction(handler *updateUserProfile.Handler) *UpdateUserProfileAction {
	return &UpdateUserProfileAction{
		Handler: handler,
	}
}

func (a *UpdateUserProfileAction) GetInputType() reflect.Type {
	return reflect.TypeOf(UpdateUserProfileRequest{})
}

type UpdateUserProfileRequest struct {
	FirstName *string `json:"firstName,omitempty"  example:"Пользователь"`
	SecondName  *string `json:"secondName,omitempty"  example:"Пользователь"`
	Address *string `json:"address,omitempty"  example:"Пользователь"`
	BirthDate *string `json:"birthDate,omitempty" example:"21.03.2026"`
	IsActive  *bool   `json:"isActive,omitempty" example:"false"`
	Email     *string `json:"email,omitempty" example:"user@example.com"`
}

func (a *UpdateUserProfileAction) Action(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromRequest(r)
	if claims == nil || claims.Phone == "" {
		a.jsonError(w, "Не удалось получить данные пользователя из токена", http.StatusUnauthorized)
		return
	}

	var req UpdateUserProfileRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		a.jsonError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	var fullName *string
	firstName := ""
	secondName := ""
	if req.FirstName != nil {
		firstName = strings.TrimSpace(*req.FirstName)
	}
	if req.SecondName != nil {
		secondName = strings.TrimSpace(*req.SecondName)
	}
	combined := strings.TrimSpace(firstName + " " + secondName)
	if combined != "" {
		fullName = &combined
	}

	cmd := updateUserProfile.Command{
		Phone:     claims.Phone,
		FullName:  fullName,
		IsActive:  req.IsActive,
		Email:     req.Email,
		BirthDate: nil,
	}

	if req.BirthDate != nil {
		raw := strings.TrimSpace(*req.BirthDate)
		if raw != "" {
			parsed, err := time.Parse("02.01.2006", raw)
			if err != nil {
				a.jsonError(w, "Некорректный формат даты рождения", http.StatusBadRequest)
				return
			}
			cmd.BirthDate = &parsed
		}
	}

	if err := a.Handler.Handle(r.Context(), cmd); err != nil {
		a.jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, UpdateUserProfileSuccessResponse{IsSuccess: true}, http.StatusOK)
}

func (a *UpdateUserProfileAction) jsonError(w http.ResponseWriter, msg string, status int) {
	utils.JSONResponse(w, UpdateUserProfileErrorResponse{
		IsSuccess: false,
		Error:     msg,
	}, status)
}

type UpdateUserProfileSuccessResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

type UpdateUserProfileErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}
