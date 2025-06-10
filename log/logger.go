package log

import "go.uber.org/zap"

var (
	globalLogger        *zap.Logger
	globalSugaredLogger *zap.SugaredLogger
)

type Config struct {
	Level      string
	File       string
	Encoder    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type logger struct {
	config Config
}

func InitLogger(opts ...Option) {
	l := &logger{
		config: Config{
			Level:      "info",
			File:       "",
			Encoder:    "console",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   false,
		},
	}
	for _, opt := range opts {
		opt(l)
	}

	globalLogger = initZapLogger(l.config)
	globalSugaredLogger = globalLogger.Sugar()

	defer globalLogger.Sync()
}

func GetLogger() *zap.Logger {
	return globalLogger
}

func GetSugaredLogger() *zap.SugaredLogger {
	return globalSugaredLogger
}

func Debug(format string, args ...any) {
	globalSugaredLogger.Debugf(format, args...)
}

func Info(format string, args ...any) {
	globalSugaredLogger.Infof(format, args...)
}

func Warn(format string, args ...any) {
	globalSugaredLogger.Warnf(format, args...)
}

func Error(format string, args ...any) {
	globalSugaredLogger.Errorf(format, args...)
}

func Fatal(format string, args ...any) {
	globalSugaredLogger.Fatalf(format, args...)
}
