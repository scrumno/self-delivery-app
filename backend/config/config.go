package config

import "github.com/scrumno/scrumno-api/shared/utils"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type JWTConfig struct {
	SecretKey []byte
}

func Load() *Config {
	secret := utils.GetEnv("JWT_SECRET", "default-dev-secret-change-in-production")
	secretKey := []byte(secret)

	return &Config{
		Server: ServerConfig{
			Port: utils.GetEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:         utils.GetEnv("DATABASE_HOST", "localhost"),
			Port:         utils.GetEnv("DATABASE_PORT", "5432"),
			Username:     utils.GetEnv("DATABASE_USERNAME", ""),
			Password:     utils.GetEnv("DATABASE_PASSWORD", ""),
			DatabaseName: utils.GetEnv("DATABASE_NAME", ""),
			SSLMode:      utils.GetEnv("DATABASE_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: secretKey,
		},
	}
}
