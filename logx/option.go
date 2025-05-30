package logx

type Option func(*logger)

func WithLevel(level string) Option {
	return func(l *logger) {
		l.config.Level = level
	}
}

func WithFile(File string) Option {
	return func(l *logger) {
		l.config.File = File
	}
}
