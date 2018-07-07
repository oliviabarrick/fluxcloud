package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// The default configuration implementation simply fetches a setting from an environment
// variable.
type DefaultConfig struct {
}

func (d *DefaultConfig) Optional(key string, defaultValue string) string {
	value := os.Getenv(strings.ToUpper(key))
	if value == "" {
		return defaultValue
	}

	return value
}

func (d *DefaultConfig) Required(key string) (string, error) {
	key = strings.ToUpper(key)

	value := os.Getenv(key)
	if value == "" {
		return "", errors.New(fmt.Sprintf("Required setting %s not set", key))
	}

	return value, nil
}
