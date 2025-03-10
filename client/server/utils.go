package server

import "strings"

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
