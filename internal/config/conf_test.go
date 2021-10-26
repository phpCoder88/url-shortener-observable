package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedServerConf = ServerConfig{
	Port:            8000,
	IdleTimeout:     2 * time.Minute,
	ReadTimeout:     5 * time.Second,
	WriteTimeout:    5 * time.Second,
	ShutdownTimeout: 15 * time.Second,
}

func TestParseServerConfig(t *testing.T) {
	conf, err := parseServerConfig()

	assert.NoError(t, err)
	assert.NotEmpty(t, conf)
}

func TestParseDBConfig(t *testing.T) {
	conf, err := parseDBConfig()

	assert.NoError(t, err)
	assert.NotEmpty(t, conf)
}

func TestGetConfig(t *testing.T) {
	conf, err := GetConfig()

	assert.NoError(t, err)
	assert.NotEmpty(t, *conf.Server)
	assert.NotEmpty(t, *conf.DB)
}

func TestGetConfigWithPort(t *testing.T) {
	err := os.Setenv("PORT", "8080")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = os.Unsetenv("PORT")
		if err != nil {
			t.Fatal(err)
		}
	}()

	conf, err := GetConfig()

	expectedServerConf.Port = 8080
	defer func() {
		expectedServerConf.Port = 8000
	}()

	assert.NoError(t, err)
	assert.NotEmpty(t, *conf.Server)
	assert.NotEmpty(t, *conf.DB)
}

func TestGetConfigWithWrongPort(t *testing.T) {
	err := os.Setenv("PORT", "nine thousand")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = os.Unsetenv("PORT")
		if err != nil {
			t.Fatal(err)
		}
	}()

	conf, err := GetConfig()

	expectedServerConf.Port = 8080
	defer func() {
		expectedServerConf.Port = 8000
	}()

	assert.Error(t, err)
	assert.Nil(t, conf)
}
