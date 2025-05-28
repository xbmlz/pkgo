package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func New(configFile string) (c *Config, err error) {
	var stat os.FileInfo

	stat, err = os.Stat(configFile)
	if err != nil {
		return
	}

	if !stat.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", configFile)
	}

	p := viper.New()
	p.SetConfigFile(configFile)

	err = p.ReadInConfig()
	if err != nil {
		return
	}

	return &Config{p}, nil
}
