package check_status_connect_db

import "github.com/scrumno/scrumno-api/internal/health/entity/status"

func Fetcher() *Query {
	s := status.Check()

	if s.IsConnected {
		return &Query{
			Status: true,
		}
	}

	return &Query{
		Status: false,
	}
}
