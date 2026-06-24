package auth

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	checkOntimeCode "github.com/scrumno/scrumno-api/internal/authorize/command/check-ontime-code"
	createAuthorizeTokens "github.com/scrumno/scrumno-api/internal/authorize/command/create-authorize-tokens"
	createUser "github.com/scrumno/scrumno-api/internal/authorize/command/create-user"
	findUserByPhone "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
)

type RegistrationAction struct {
	FindUserByPhoneFetcher       *findUserByPhone.Fetcher
	CheckOntimeCodeHandler       *checkOntimeCode.Handler
	CreateUserHandler            *createUser.Handler
	CreateAuthorizeTokensHandler *createAuthorizeTokens.Handler
}

func NewRegistrationAction(
	findUserByPhoneFetcher *findUserByPhone.Fetcher,
	checkOntimeCodeHandler *checkOntimeCode.Handler,
	createUserHandler *createUser.Handler,
	createAuthorizeTokensHandler *createAuthorizeTokens.Handler,
) *RegistrationAction {
	return &RegistrationAction{
		FindUserByPhoneFetcher:       findUserByPhoneFetcher,
		CheckOntimeCodeHandler:       checkOntimeCodeHandler,
		CreateUserHandler:            createUserHandler,
		CreateAuthorizeTokensHandler: createAuthorizeTokensHandler,
	}
}

func (a *RegistrationAction) GetInputType() reflect.Type {
	return reflect.TypeOf(RegistrationRequest{})
}

type RegistrationRequest struct {
	Phone string `json:"phone" example:"79099000000"`
	Code  string `json:"code" example:"1234"`
}

func (a *RegistrationAction) Action(w http.ResponseWriter, r *http.Request) {

	var req RegistrationRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	user, err := a.FindUserByPhoneFetcher.Fetch(r.Context(), req.Phone)
	if err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if user != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     "Не удалось создать пользователя",
			Code:      "USER_EXIST",
		}, http.StatusBadRequest)
		return
	}

	// cmd := checkOntimeCode.Command{
	// 	Phone:    req.Phone,
	// 	Code:     req.Code,
	// 	CodeType: codes.RegisterType,
	// }

	// err = a.CheckOntimeCodeHandler.Handle(r.Context(), cmd)
	// if err != nil {
	// 	utils.JSONResponse(w, RegistrationErrorResponse{
	// 		IsSuccess: false,
	// 		Error:     err.Error(),
	// 	}, http.StatusBadRequest)
	// 	return
	// }

	createdUser, err := a.CreateUserHandler.Handle(
		r.Context(),
		createUser.Command{
			Phone: req.Phone,
		},
	)
	if err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, expiresIn, err := a.CreateAuthorizeTokensHandler.Handle(
		r.Context(),
		createAuthorizeTokens.Command{
			Phone:               createdUser.Phone,
			UserID:              createdUser.ID,
			SessionID:           "",
			RevokePreviousToken: false,
		},
	)
	if err != nil {
		utils.JSONResponse(w, RegistrationErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, RegistrationResponse{
		IsSuccess:    true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, http.StatusOK)
}

type RegistrationResponse struct {
	IsSuccess    bool   `json:"isSuccess"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type RegistrationErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
}
