package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LoggerKey = "logger"

type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func InitLogger(level, format string) error {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) SInfo(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) SDebug(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) SError(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *Logger) SWarn(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) SFatal(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func Sync() {
	_ = zap.L().Sync()
}

func GetLogger(ctx context.Context) *Logger {
	if l, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		return &Logger{
			logger: l,
			sugar:  l.Sugar(),
		}
	}
	return &Logger{
		logger: zap.L(),
		sugar:  zap.L().Sugar(),
	}
}
