package status

import "github.com/scrumno/scrumno-api/config"

func Check() *Status {
	sqlDB, err := config.DB.DB()
	if err != nil {
		return &Status{IsConnected: false}
	}

	if err := sqlDB.Ping(); err != nil {
		return &Status{IsConnected: false}
	}

	return &Status{IsConnected: true}
}
