package config

import (
	"errors"
	"fmt"
	"strings"
)

type FakeConfig struct {
	settings map[string]string
}

func NewFakeConfig() *FakeConfig {
	return &FakeConfig{
		settings: map[string]string{},
	}
}

func (t *FakeConfig) Optional(key string, defaultValue string) string {
	value := t.settings[strings.ToUpper(key)]
	if value == "" {
		return defaultValue
	}

	return value
}

func (t *FakeConfig) Required(key string) (string, error) {
	value := t.settings[strings.ToUpper(key)]
	if value == "" {
		return "", errors.New(fmt.Sprintf("Required setting %s not set", key))
	}

	return value, nil
}

func (t *FakeConfig) Set(key string, value string) {
	t.settings[strings.ToUpper(key)] = value
}
