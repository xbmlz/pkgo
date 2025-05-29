package log

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initZapLogger(cfg Config) *zap.Logger {
	cores := []zapcore.Core{
		createConsoleZapCore(cfg.Level, cfg.Encoder),
	}
	if cfg.FileName != "" {
		if _, err := os.Stat(cfg.FileName); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(cfg.FileName), os.ModePerm)
		}
		cores = append(cores, createFileZapCore(cfg))
	}
	zapOpts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}
	if cfg.Level == "debug" {
		zapOpts = append(zapOpts, zap.Development(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	logger := zap.New(zapcore.NewTee(cores...), zapOpts...)

	defer logger.Sync()

	return logger
}

func createConsoleZapCore(level, encoder string) zapcore.Core {
	return zapcore.NewCore(
		createEncoderConfig(encoder),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zap.NewAtomicLevelAt(parseLevel(level)),
	)
}

func createFileZapCore(cfg Config) zapcore.Core {
	hook := &lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.NewCore(
		createEncoderConfig(cfg.Encoder),
		zapcore.AddSync(hook),
		zap.NewAtomicLevelAt(parseLevel(cfg.Level)),
	)
}

func createEncoderConfig(encoder string) zapcore.Encoder {
	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	consoleEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	consoleEncoderConfig.EncodeName = zapcore.FullNameEncoder
	consoleEncoderConfig.ConsoleSeparator = "\t"
	if encoder == "json" {
		return zapcore.NewJSONEncoder(consoleEncoderConfig)
	}
	return zapcore.NewConsoleEncoder(consoleEncoderConfig)
}

func parseLevel(level string) zapcore.Level {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return zapcore.InfoLevel
	}
	return l
}
