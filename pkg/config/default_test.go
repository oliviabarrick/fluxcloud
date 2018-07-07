package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDefaultConfigImplementsConfig(t *testing.T) {
	_ = Config(&DefaultConfig{})

}

func TestDefaultConfigOptional(t *testing.T) {
	config := DefaultConfig{}

	os.Setenv("GIT_URL", "github.com")
	assert.Equal(t, "github.com", config.Optional("git_url", "hello"))

	os.Setenv("GIT_URL", "")
	assert.Equal(t, "hello", config.Optional("git_url", "hello"))
}

func TestDefaultConfigRequired(t *testing.T) {
	config := DefaultConfig{}

	os.Setenv("GIT_URL", "github.com")
	gitUrl, err := config.Required("git_url")
	assert.Nil(t, err)
	assert.Equal(t, "github.com", gitUrl)

	os.Setenv("GIT_URL", "")
	gitUrl, err = config.Required("git_url")
	assert.NotNil(t, err)
	assert.Equal(t, "", gitUrl)
}
