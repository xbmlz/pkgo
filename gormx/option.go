package gormx

type Option func(*Config)

func WithMaxIdleConns(n int) Option {
	return func(c *Config) {
		c.MaxIdleConns = n
	}
}

func WithMaxOpenConns(n int) Option {
	return func(c *Config) {
		c.MaxOpenConns = n
	}
}

func WithConnMaxLifeTime(t int) Option {
	return func(c *Config) {
		c.ConnMaxLifeTime = t
	}
}
