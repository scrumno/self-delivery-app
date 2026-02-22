package user

import "github.com/scrumno/scrumno-api/shared/base"

type userRepository interface {
	base.BaseRepository[User]
}
