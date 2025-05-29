package log

type Option func(*logger)

func WithLevel(level string) Option {
	return func(l *logger) {
		l.config.Level = level
	}
}

func WithFileName(filename string) Option {
	return func(l *logger) {
		l.config.FileName = filename
	}
}
