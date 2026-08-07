package enval

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func Get(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func GetInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func GetDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}

func GetAsBool(name string, defaultVal bool) bool {
	valStr := Get(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func GetAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := Get(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}

func GetAsByte(name string, defaultVal []byte) []byte {
	valStr := Get(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := []byte(strings.Replace(valStr, `\n`, "\n", -1))

	return val
}
