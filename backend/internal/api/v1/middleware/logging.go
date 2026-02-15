package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// LogEntry представляет структуру лога со всеми деталями
type LogEntry struct {
	Timestamp   string              `json:"timestamp"`
	Method      string              `json:"method"`
	Duration    string              `json:"duration"`
	URL         string              `json:"url"`
	Path        string              `json:"path"`
	Query       string              `json:"query,omitempty"`
	QueryParams map[string][]string `json:"query_params,omitempty"`
	Headers     http.Header         `json:"headers"`
	Body        string              `json:"body,omitempty"`
	Status      int                 `json:"status"`
	RemoteAddr  string              `json:"remote_addr"`
	UserAgent   string              `json:"user_agent"`
	ContentType string              `json:"content_type,omitempty"`
	RequestID   string              `json:"request_id,omitempty"`
}

type bodyReader struct {
	io.ReadCloser
	body *bytes.Buffer
}

func (br *bodyReader) Read(p []byte) (int, error) {
	n, err := br.ReadCloser.Read(p)
	if n > 0 {
		br.body.Write(p[:n]) // Сохраняем прочитанные данные
	}
	return n, err
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		var bodyBytes []byte
		var bodyBuffer *bytes.Buffer

		if r.Body != nil {
			bodyBuffer = bytes.NewBuffer(nil)
			r.Body = &bodyReader{
				ReadCloser: r.Body,
				body:       bodyBuffer,
			}

			defer func() {
				if bodyBuffer.Len() > 0 {
					bodyBytes = bodyBuffer.Bytes()
				}
			}()
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		path := r.URL.Path
		query := r.URL.RawQuery

		queryParams := make(map[string][]string)
		for key, values := range r.URL.Query() {
			queryParams[key] = values
		}

		headers := make(http.Header)
		for key, values := range r.Header {
			if strings.ToLower(key) == "authorization" ||
				strings.ToLower(key) == "cookie" ||
				strings.ToLower(key) == "set-cookie" {
				headers[key] = []string{"[REDACTED]"}
			} else {
				headers[key] = values
			}
		}

		bodyString := ""
		if len(bodyBytes) > 0 {
			if json.Valid(bodyBytes) {
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
					bodyString = prettyJSON.String()
				} else {
					bodyString = string(bodyBytes)
				}
			} else {
				maxBodyLen := 1000
				if len(bodyBytes) > maxBodyLen {
					bodyString = string(bodyBytes[:maxBodyLen]) + "... [truncated]"
				} else {
					bodyString = string(bodyBytes)
				}
			}
		}

		entry := LogEntry{
			Timestamp:   start.Format(time.RFC3339Nano),
			Method:      r.Method,
			Duration:    duration.String(),
			URL:         r.URL.String(),
			Path:        path,
			Query:       query,
			QueryParams: queryParams,
			Headers:     headers,
			Body:        bodyString,
			Status:      rw.statusCode,
			RemoteAddr:  getRealIP(r),
			UserAgent:   r.UserAgent(),
			ContentType: r.Header.Get("Content-Type"),
			RequestID:   getRequestID(r),
		}

		jsonData, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
			return
		}

		fmt.Fprintln(os.Stdout, string(jsonData))
	})
}

func getRealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// Убираем порт из RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func getRequestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	if id := r.Header.Get("X-Trace-ID"); id != "" {
		return id
	}
	return ""
}
