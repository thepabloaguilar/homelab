package main

import (
	"github.com/spf13/cobra"

	"github.com/thepabloaguilar/homelab/internal/config"

	"github.com/thepabloaguilar/homelab/internal/ssh_operations"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Set of commands to manage the server",
}

var serverSetupCmd = &cobra.Command{
	Use:   "hello",
	Short: "Hello to the servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cfg := config.FromContext(ctx)

		hello := ssh_operations.Hello()
		for _, server := range cfg.Servers {
			if err := hello(ctx, server); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	serverCmd.AddCommand(serverSetupCmd)
}
