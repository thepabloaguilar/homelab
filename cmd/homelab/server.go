package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thepabloaguilar/homelab/internal/config"
	"github.com/thepabloaguilar/homelab/internal/tools"

	"github.com/thepabloaguilar/homelab/internal/ssh_operations"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Set of commands to manage the server",
}

var serverSSHSetupCmd = &cobra.Command{
	Use:       "setup-ssh [server_name]",
	Short:     "Setup the server ssh configs",
	ValidArgs: []string{"server_name"},
	Args:      cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cfg := config.FromContext(ctx)
		server, ok := cfg.Servers[args[0]]
		if !ok {
			return fmt.Errorf("server not found: %s", args[0])
		}

		sshFiles, err := cmd.Flags().GetString("ssh-files")
		if err != nil {
			return err
		}

		client, err := tools.NewSSHClient(server)
		if err != nil {
			return err
		}
		defer tools.LogCloser(ctx, client)

		err = ssh_operations.NewHardenSSHSecurity(os.DirFS(sshFiles)).Do(ctx, client)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	serverSSHSetupCmd.Flags().String(
		"ssh-files",
		"./assets/setup_ssh",
		"SSH configuration files to be copied",
	)

	serverCmd.AddCommand(serverSSHSetupCmd)
}
