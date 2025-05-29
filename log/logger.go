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

func Debug(args ...any) {
	globalSugaredLogger.Debug(args...)
}

func Debugf(format string, args ...any) {
	globalSugaredLogger.Debugf(format, args...)
}

func Info(args ...any) {
	globalSugaredLogger.Info(args...)
}

func Infof(format string, args ...any) {
	globalSugaredLogger.Infof(format, args...)
}

func Warn(args ...any) {
	globalSugaredLogger.Warn(args...)
}

func Warnf(format string, args ...any) {
	globalSugaredLogger.Warnf(format, args...)
}

func Error(args ...any) {
	globalSugaredLogger.Error(args...)
}

func Errorf(format string, args ...any) {
	globalSugaredLogger.Errorf(format, args...)
}

func Fatal(args ...any) {
	globalSugaredLogger.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	globalSugaredLogger.Fatalf(format, args...)
}
