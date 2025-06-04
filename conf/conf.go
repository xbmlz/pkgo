package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func Load(configFile string) (c *Config, err error) {
	p := viper.New()
	p.SetConfigFile(configFile)

	err = p.ReadInConfig()
	if err != nil {
		return
	}

	return &Config{p}, nil
}

func (c *Config) Parse(v any) error {
	return c.Viper.Unmarshal(v)
}

func MustLoad(configFile string) *Config {
	c, err := Load(configFile)
	if err != nil {
		panic(err)
	}
	return c
}
