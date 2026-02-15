package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/scrumno/scrumno-api/config"
	v1 "github.com/scrumno/scrumno-api/internal/api/v1"
	jwtEntity "github.com/scrumno/scrumno-api/internal/users/auth/jwt/entity"
	"github.com/scrumno/scrumno-api/internal/users/user/entity/user"
)

func main() {
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend/.env.local")
	_ = godotenv.Load("backend/.env")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()
	if err := config.Connect(cfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := config.Migrate(&user.User{}, &jwtEntity.UserToken{}); err != nil {
		logger.Error("миграция БД", "error", err)
		os.Exit(1)
	}

	defer func() {
		err := config.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	router := v1.SetupRouter(cfg)
	addr := ":" + cfg.Server.Port
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Сервер запущен", "address", addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Сервер не запустился", "error", err)
		os.Exit(1)
	}
}
