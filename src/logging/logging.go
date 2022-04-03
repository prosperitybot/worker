package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() error {
	var err error

	config := zap.NewProductionConfig()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderCfg

	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	defer log.Sync() // flushes buffer, if any

	return nil
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields)
	log.Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields)
	log.Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields)
	log.Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields)
	log.Error(msg, fields...)
}
