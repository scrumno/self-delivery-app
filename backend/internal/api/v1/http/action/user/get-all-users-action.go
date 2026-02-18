package user

import "net/http"

// Получает, берёт через fetcher query

func GetAllUserAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "list of users"}`))
}
