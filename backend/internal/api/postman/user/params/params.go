package params

type CreateUserRequest struct {
    Phone    string  `json:"phone" example:"79099009988"`
	FullName string `json:"full_name" example:"Иван Иванов"`
}

type GetUserRequest struct {}
