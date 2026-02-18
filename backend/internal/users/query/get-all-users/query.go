package get_all_users

import "github.com/scrumno/scrumno-api/internal/users/entity/user"

type Query struct {
	userData user.User
}
