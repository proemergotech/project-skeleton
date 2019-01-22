package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/di"
	"gitlab.com/proemergotech/log-go"
)

var (
	cfg = &config.Config{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: config.AppName,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		container, err := di.NewContainer(cfg)
		if err != nil {
			log.Panic(context.Background(), "Couldn't load container", "error", err)
		}

		defer func() {
			container.Close()
		}()

		errorCh := make(chan error)
		container.RestServer.Start(errorCh)

		defer func() {
			err = container.RestServer.Stop(5 * time.Second)
			if err != nil {
				err = errors.Wrap(err, "Rest server graceful shutdown failed")
				log.Panic(context.Background(), err.Error(), "error", err)
			}
			log.Info(context.Background(), "Shutdown complete")
		}()

		log.Info(context.Background(), "Rest server started")

		container.EventServer.Start()
		log.Info(context.Background(), "Event server started")

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigs:
		case err := <-errorCh:
			err = errors.Wrap(err, "Rest server fatal error")
			log.Panic(context.Background(), err.Error(), "error", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
