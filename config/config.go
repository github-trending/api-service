package config

import (
	"os"
)

var defaultConfig = map[string]string{
	"debug":      "false",
	"host":       "http://localhost:8080",
	"redis_addr": "127.0.0.1:6379",
	"redis_auth": "",
}

func Get(name string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultConfig[name]
	}

	return value
}
