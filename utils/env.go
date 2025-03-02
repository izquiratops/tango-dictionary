package utils

import (
	"os"
	"strings"
)

func ResolveBooleanFromEnv(envName string) bool {
	value := strings.ToLower(os.Getenv(envName))
	return value == "true" || value == "1" || value == "yes"
}
