package check_status_connect_db

import (
	"github.com/scrumno/scrumno-api/internal/health/entity/status"
)

type StatusDTO struct {
	IsConnected bool `json:"is_connected"`
}

type Fetcher struct {
	repository status.StatusRepositoryInterface
}

func NewFetcher(repository status.StatusRepositoryInterface) *Fetcher {
	return &Fetcher{repository: repository}
}

func (f *Fetcher) Fetch(_ Query) StatusDTO {
	err := f.repository.CheckStatus()
	return StatusDTO{IsConnected: err == nil}
}
