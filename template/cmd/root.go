package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/proemergotech/errors"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/trace/v2"
	"github.com/spf13/cobra"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/config"
	"github.com/proemergotech/project-skeleton/app/di"
	//%: ` | replace "project-skeleton" .ProjectName }}
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

		runner := newRunner()
		defer runner.stop()

		runner.start("rest server", container.RestServer.Start, container.RestServer.Stop)
		//%: {{- if .PublicRest }}
		runner.start("public rest server", container.PublicRestServer.Start, container.PublicRestServer.Stop)
		//%: {{- end }}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigs:
		case err := <-runner.errors():
			ctx := context.Background()

			correlationID := errors.Field(err, trace.CorrelationIDField)
			if correlationID != nil {
				ctx = trace.WithCorrelation(ctx, &trace.Correlation{
					CorrelationID: correlationID.(string),
				})
			}

			log.Panic(ctx, err.Error(), "error", err)
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
