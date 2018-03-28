package config

import (
  "os"
)

var defaultConfig = map[string]string{
	"debug": "false",
	"host": "http://localhost:8080",
}

func Get(name string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultConfig[name]
	}

	return value
}
