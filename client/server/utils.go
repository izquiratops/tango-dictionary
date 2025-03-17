package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (s *Server) logRequest(r *http.Request, statusCode int, duration time.Duration) {
	ip := r.RemoteAddr
	query := r.URL.Query().Get("query")

	// Format: [TIMESTAMP] METHOD PATH STATUS IP DURATION QUERY
	logMsg := fmt.Sprintf("[%s] %s %s %d %s %v %s",
		time.Now().Format("2006-01-02 15:04:05"),
		r.Method,
		r.URL.Path,
		statusCode,
		ip,
		duration,
		query)

	fmt.Println(logMsg)
}

func getContentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".ttf"):
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}
