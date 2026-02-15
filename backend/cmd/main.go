package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Response Структура ответа JSON
type Response struct {
	Message string `json:"message"`
}

func main() {
	// 1. Инициализация логгера (JSON формат удобен для продакшена/CloudWatch/ELK)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// 2. Создаем роутер Gorilla
	r := mux.NewRouter()

	// 3. Регистрируем хендлер
	r.HandleFunc("/api/v1/health/simple-check-status", func(w http.ResponseWriter, r *http.Request) {
		// Логируем запрос
		logger.Info("Health check requested",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Отправляем JSON
		json.NewEncoder(w).Encode(Response{
			Message: "hello world",
		})
	}).Methods("GET")

	// 4. Настройки сервера
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Server started", "address", ":8080", "go_version", "1.25.7")

	// 5. Запуск
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
