package main

import (
	"github.com/spf13/cobra"

	"github.com/thepabloaguilar/homelab/internal/config"
	"github.com/thepabloaguilar/homelab/internal/tools"
)

var rootCmd = &cobra.Command{
	Use:   "homelab",
	Short: "CLI to help manage homelab",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		logger := tools.NewLogger()
		cmd.SetContext(
			tools.LoggerToContext(
				config.ToContext(cmd.Context(), cfg), logger,
			),
		)

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		_ = tools.LoggerFromContext(cmd.Context()).Sync()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
