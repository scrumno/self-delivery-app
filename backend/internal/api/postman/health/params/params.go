package params

type CreateUserRequest struct {
    Phone    string  `json:"phone" example:"79099000000`
	FullName string `json:"full_name" example:"Иван Аресньев"`
}

type GetUserRequest struct {}
