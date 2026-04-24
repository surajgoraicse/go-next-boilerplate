package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	level, err := parseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    maxOrDefault(cfg.LogMaxSizeMB, 100),
		MaxBackups: maxOrDefault(cfg.LogMaxBackups, 5),
		MaxAge:     maxOrDefault(cfg.LogMaxAgeDays, 7),
		Compress:   true,
	}

	if err := ensureLogDir(cfg.LogFile); err != nil {
		return nil, err
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), level),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.ErrorOutput(zapcore.AddSync(os.Stderr)),
	).With(
		zap.String("service", cfg.AppName),
		zap.String("env", string(cfg.AppEnv)),
	)

	return logger, nil
}

func parseLevel(v string) (zapcore.Level, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "info":
		return zapcore.InfoLevel, nil
	case "debug":
		return zapcore.DebugLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s", v)
	}
}

func ensureLogDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o755)
}

func maxOrDefault(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}
