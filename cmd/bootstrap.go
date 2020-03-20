package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
)

var bootstrapCmd = &cobra.Command{
	Use: "bootstrap",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.BootstrapConfig{}
		initConfig(cfg)

		// do something
	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}
