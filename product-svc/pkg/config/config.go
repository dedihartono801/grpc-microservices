package config

import (
	"os"
)

// GetEnv
func GetEnv(key string) string {
	value := os.Getenv(key)
	return value
}
