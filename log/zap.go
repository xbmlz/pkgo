package log

import (
	"log"
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
	if cfg.File != "" {
		cores = append(cores, createFileZapCore(cfg))
	}

	return zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),                   // 记录调用位置
		zap.AddCallerSkip(1),              // 包装层数
		zap.AddStacktrace(zap.ErrorLevel), // Error及以上记录堆栈
	)

}

func createConsoleZapCore(level, encoder string) zapcore.Core {
	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	consoleEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	consoleEncoderConfig.EncodeName = zapcore.FullNameEncoder
	consoleEncoderConfig.ConsoleSeparator = "\t"
	return zapcore.NewCore(
		createEncoderConfig(encoder, consoleEncoderConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zap.NewAtomicLevelAt(ParseLevel(level)),
	)
}

func createFileZapCore(cfg Config) zapcore.Core {
	if err := os.MkdirAll(filepath.Dir(cfg.File), os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	fileEncoderConfig.EncodeName = zapcore.FullNameEncoder
	fileEncoderConfig.ConsoleSeparator = "\t"
	hook := &lumberjack.Logger{
		Filename:   cfg.File,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.NewCore(
		createEncoderConfig(cfg.Encoder, fileEncoderConfig),
		zapcore.AddSync(hook),
		zap.NewAtomicLevelAt(ParseLevel(cfg.Level)),
	)
}

func createEncoderConfig(encoder string, consoleEncoderConfig zapcore.EncoderConfig) zapcore.Encoder {
	if encoder == "json" {
		return zapcore.NewJSONEncoder(consoleEncoderConfig)
	}
	return zapcore.NewConsoleEncoder(consoleEncoderConfig)
}

func ParseLevel(level string) zapcore.Level {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return zapcore.InfoLevel
	}
	return l
}
