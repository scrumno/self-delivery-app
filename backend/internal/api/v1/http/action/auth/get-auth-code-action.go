package auth

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	createAuthorizeCode "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-code"
	codes "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"
	getSmsCode "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code"
	getSmsCodeSendAvailable "github.com/scrumno/scrumno-api/internal/authorize/query/get-sms-code-send-available"
	findUserByPhone "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
)

type AuthCodeAction struct {
	GetSmsCodeSendAvailableFetcher *getSmsCodeSendAvailable.Fetcher
	GetSmsCodeFetcher              *getSmsCode.Fetcher
	CreateAuthorizeCodeHandler     *createAuthorizeCode.Handler
	FindUserByPhoneFetcher         *findUserByPhone.Fetcher
}

func NewAuthCodeAction(
	getSmsCodeSendAvailableFetcher *getSmsCodeSendAvailable.Fetcher,
	getSmsCodeFetcher *getSmsCode.Fetcher,
	createAuthorizeCodeHandler *createAuthorizeCode.Handler,
	findUserByPhoneFetcher *findUserByPhone.Fetcher,
) *AuthCodeAction {
	return &AuthCodeAction{
		GetSmsCodeSendAvailableFetcher: getSmsCodeSendAvailableFetcher,
		GetSmsCodeFetcher:              getSmsCodeFetcher,
		CreateAuthorizeCodeHandler:     createAuthorizeCodeHandler,
		FindUserByPhoneFetcher:         findUserByPhoneFetcher,
	}
}

func (a *AuthCodeAction) GetInputType() reflect.Type {
	return reflect.TypeOf(GetAuthCodeRequest{})
}

type GetAuthCodeRequest struct {
	Phone    string          `json:"phone" example:"79090000000"`
	CodeType codes.CodesType `json:"codeType" example:"authorize"`
}

func (a *AuthCodeAction) Action(w http.ResponseWriter, r *http.Request) {

	var req GetAuthCodeRequest
	err := utils.DecodeJSONBody(r, &req)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Неверный формат запроса",
		}, http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Укажите номер телефона",
		}, http.StatusBadRequest)
		return
	}

	if req.CodeType != codes.AuthType && req.CodeType != codes.RegisterType {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Недопустимый тип кода",
		}, http.StatusBadRequest)
		return
	}

	
	user, err := a.FindUserByPhoneFetcher.Fetch(r.Context(), req.Phone)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if user == nil && req.CodeType == codes.AuthType {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Пользователь не найден",
			Code:      "USER_NOT_FOUND",
		}, http.StatusBadRequest)
		return
	}

	if user != nil && req.CodeType == codes.RegisterType {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     "Не удалось получить код авторизации",
			Code:      "USER_EXIST",
		}, http.StatusBadRequest)
		return
	}

	_, err = a.GetSmsCodeSendAvailableFetcher.Fetch(r.Context(), req.Phone)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	authorizeCode, err := a.CreateAuthorizeCodeHandler.Handle(
		r.Context(),
		req.Phone,
		req.CodeType,
	)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("Ваш код авторизации: %s", authorizeCode.Code)

	_, err = a.GetSmsCodeFetcher.Fetch(r.Context(), req.Phone, message)
	if err != nil {
		utils.JSONResponse(w, AuthCodeErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, AuthCodeResponse{
		IsSuccess: true,
	}, http.StatusOK)
}

type AuthCodeResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

type AuthCodeErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
}
