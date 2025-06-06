package conf

import (
	"github.com/spf13/viper"
)

type option struct {
	configFile string
}

type OptionFunc func(*option)

func WithConfigFile(configFile string) OptionFunc {
	return func(o *option) {
		o.configFile = configFile
	}
}

func Load(v any, opts ...OptionFunc) error {
	opt := option{
		configFile: "config.yaml",
	}

	for _, o := range opts {
		o(&opt)
	}

	p := viper.New()
	p.SetConfigFile(opt.configFile)

	if err := p.ReadInConfig(); err != nil {
		return err
	}

	if err := p.Unmarshal(v); err != nil {
		return err
	}

	return nil
}
