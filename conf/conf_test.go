package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Server struct {
		Host string `yaml:"host" json:"host"`
		Port int    `yaml:"port" json:"port"`
	} `yaml:"server" json:"server"`
	Log struct {
		Level string `yaml:"level" json:"level"`
	} `yaml:"log" json:"log"`
}

func TestLoad(t *testing.T) {
	var c TestConfig

	config, err := Load("./testdata/config.yaml")

	assert.NoError(t, err)
	assert.NotNil(t, config)

	err = config.Unmarshal(&c)
	assert.NoError(t, err)

	assert.Equal(t, c.Server.Host, "0.0.0.0")
	assert.Equal(t, c.Server.Port, 8080)
	assert.Equal(t, c.Log.Level, "debug")
}
