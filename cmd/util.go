package cmd

import (
	"os"
	"strings"
)

func env(key string, defaults ...string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	if defaults != nil {
		return strings.Join(defaults, "")
	}
	return ""
}

func envBool(key string) bool {
	return strings.ToLower(env(key)) == "true" || env(key) == "1"
}
