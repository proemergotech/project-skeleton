package main

import (
	"context"
	"fmt"
	"os"

	"github.com/proemergotech/errors"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/log/v3/zaplog"
	"github.com/proemergotech/trace/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/config"
	"github.com/proemergotech/project-skeleton/cmd"
	//%: ` | replace "project-skeleton" .ProjectName }}
)

func main() {
	if err := zap.RegisterEncoder(zaplog.EncoderType, zaplog.NewEncoder([]string{
		trace.CorrelationIDField,
		trace.WorkflowIDField,
		log.AppName,
		log.AppVersion,
	})); err != nil {
		panic(fmt.Sprintf("Couldn't create logger, error: %v", err))
	}

	zapConf := zap.NewProductionConfig()
	zapConf.Encoding = zaplog.EncoderType

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	zapLogLevel := new(zapcore.Level)
	if err := zapLogLevel.Set(logLevel); err != nil {
		panic(fmt.Sprintf("Invalid log level: %s", logLevel))
	}
	zapConf.Level = zap.NewAtomicLevelAt(*zapLogLevel)

	zapLog, err := zapConf.Build()
	if err != nil {
		panic(fmt.Sprintf("Couldn't create logger, error: %v", err))
	}
	zapLog = zapLog.With(
		zap.String(log.AppName, config.AppName),
		zap.String(log.AppVersion, config.AppVersion),
	)
	log.SetGlobalLogger(zaplog.NewLogger(zapLog, trace.Mapper()))

	defer func() {
		if err := recover(); err != nil {
			log.Error(context.Background(), "Service panicked", "error", errors.Errorf("%+v", err))
			os.Exit(1)
		}
	}()

	cmd.Execute()
}
