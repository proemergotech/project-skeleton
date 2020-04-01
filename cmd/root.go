package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/di"
	"gitlab.com/proemergotech/errors"
	log "gitlab.com/proemergotech/log-go/v3"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: config.AppName,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{}
		initConfig(cfg)

		container, err := di.NewContainer(cfg)
		if err != nil {
			log.Panic(context.Background(), "Couldn't load container", "error", err)
		}
		defer container.Close()

		errorCh := make(chan error)
		container.RestServer.Start(errorCh)

		defer func() {
			if err := container.RestServer.Stop(5 * time.Second); err != nil {
				err = errors.Wrap(err, "Rest server graceful shutdown failed")
				log.Panic(context.Background(), err.Error(), "error", err)
			}
			log.Info(context.Background(), "Shutdown complete")
		}()

		log.Info(context.Background(), "Rest server started")

		publicRestErrorCh := make(chan error)
		container.PublicRestServer.Start(publicRestErrorCh)

		defer func() {
			if err := container.PublicRestServer.Stop(5 * time.Second); err != nil {
				err = errors.Wrap(err, "Public rest server graceful shutdown failed")
				log.Panic(context.Background(), err.Error(), "error", err)
			}
			log.Info(context.Background(), "Public rest server graceful shutdown complete")
		}()

		log.Info(context.Background(), "Public rest server started")

		if err := container.EventServer.Start(); err != nil {
			err = errors.Wrap(err, "Failed starting event server")
			log.Panic(context.Background(), err.Error(), "error", err)
		}
		log.Info(context.Background(), "Event server started")

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigs:
		case err := <-errorCh:
			err = errors.Wrap(err, "Rest server fatal error")
			log.Panic(context.Background(), err.Error(), "error", err)
		case err := <-publicRestErrorCh:
			err = errors.Wrap(err, "Public rest server fatal error")
			log.Panic(context.Background(), err.Error(), "error", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic(context.Background(), err.Error(), "error", err)
	}
}
