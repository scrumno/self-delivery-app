package utils

import (
	"os"
	"regexp"
)

// NormalizePhone оставляет только цифры номера телефона (username).
func NormalizePhone(phone string) string {
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(phone, "")
}

// ValidatePhone проверяет формат номера: строго 79996663355 (11 цифр, первая 7).
func ValidatePhone(phone string) bool {
	normalized := NormalizePhone(phone)
	return len(normalized) == 11 && normalized[0] == '7'
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
