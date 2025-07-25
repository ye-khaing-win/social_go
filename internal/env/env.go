package env

import (
	"os"
	"strconv"
)

func GetStr(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}
