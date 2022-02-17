package config

import (
	"os"
	"strconv"
)

// GetEnvStr
func GetEnvStr(key string, default_val string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return default_val
}

// GetEnvInt
func GetEnvInt(key string, default_val int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return default_val
}

// GetEnvBool
func GetEnvBool(key string, default_val bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return default_val
}

// GetEnvFloat
func GetEnvFloat(key string, default_val float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return default_val
}
