package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Error encoding response:", err)
		return
	}

	return
}

type UpdateRequest[T any] struct {
	ID     uuid.UUID
	Fields T
	Set    map[string]struct{} // какие поля пришли
}

func DecodeUpdateRequest[T any](r *http.Request) (UpdateRequest[T], error) {
	var result UpdateRequest[T]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return result, fmt.Errorf("read body: %w", err)
	}
	defer r.Body.Close()

	// Парсим верхний уровень чтобы достать ID
	var envelope struct {
		ID   uuid.UUID       `json:"id"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return result, fmt.Errorf("unmarshal envelope: %w", err)
	}

	result.ID = envelope.ID

	if result.ID == uuid.Nil {
		return result, fmt.Errorf("id is required")
	}

	// Парсим data в нужный тип
	if err := json.Unmarshal(envelope.Data, &result.Fields); err != nil {
		return result, fmt.Errorf("unmarshal fields: %w", err)
	}

	// Определяем какие поля реально пришли
	var rawFields map[string]json.RawMessage
	if err := json.Unmarshal(envelope.Data, &rawFields); err != nil {
		return result, fmt.Errorf("unmarshal raw fields: %w", err)
	}

	result.Set = make(map[string]struct{}, len(rawFields))
	for key := range rawFields {
		result.Set[key] = struct{}{}
	}

	return result, nil
}

// IsSet проверяет, пришло ли поле в запросе
func (u *UpdateRequest[T]) IsSet(field string) bool {
	_, ok := u.Set[field]
	return ok
}

func BuildUpdateMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)

		// пропускаем nil указатели
		if val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]
		result[jsonKey] = val.Interface()
	}

	return result
}
