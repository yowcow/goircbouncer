package config

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_returns_error_on_no_file(t *testing.T) {
	cfg, err := Load("hogefuga")

	assert.Nil(t, cfg)
	assert.Regexp(t, regexp.MustCompile("no such file or directory"), err.Error())
}

func TestLoad_returns_error_on_not_file(t *testing.T) {
	cfg, err := Load("./")

	assert.Nil(t, cfg)
	assert.Regexp(t, regexp.MustCompile("is a directory"), err.Error())
}

func TestLoad_returns_error_on_invalid_yaml(t *testing.T) {
	cfg, err := Load("config-invalid.yml")

	assert.Nil(t, cfg)
	assert.Regexp(t, regexp.MustCompile("unmarshal errors"), err.Error())
}

func TestLoad_returns_config_on_valid_yaml(t *testing.T) {
	cfg, err := Load("config-valid.yml")

	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "hoge", cfg.Server.Host)
	assert.Equal(t, ":1234", cfg.Server.Addr)
	assert.Equal(t, true, cfg.Server.Secure)
	assert.Equal(t, "", cfg.Server.Password)
	assert.Equal(t, 2, len(cfg.Server.Channels))
	assert.Equal(t, "#foo", cfg.Server.Channels[0])
	assert.Equal(t, "#bar", cfg.Server.Channels[1])

	assert.Equal(t, ":2345", cfg.Client.Addr)
	assert.Equal(t, "hogefuga", cfg.Client.Password)
}
