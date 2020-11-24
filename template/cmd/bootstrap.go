//%: {{ if .Bootstrap }}
package cmd

import (
	"github.com/spf13/cobra"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/config"
	//%: ` | replace "project-skeleton" .ProjectName }}
)

// todo: remove OR update
//  Example bootstrap subcommand
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

//%: {{ end }}
