package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/cmd"
	log "gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/log-go/zaplog"
	trace "gitlab.com/proemergotech/trace-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	err := zap.RegisterEncoder(zaplog.EncoderType, zaplog.NewEncoder([]string{
		trace.CorrelationIDField,
		trace.WorkflowIDField,
		log.AppName,
		log.AppVersion,
	}))
	if err != nil {
		panic(fmt.Sprintf("Couldn't create logger, error: %v", err))
	}

	zapConf := zap.NewProductionConfig()
	zapConf.Encoding = zaplog.EncoderType

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	zapLogLevel := new(zapcore.Level)
	err = zapLogLevel.Set(logLevel)
	if err != nil {
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
