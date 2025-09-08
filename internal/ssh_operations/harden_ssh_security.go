package ssh_operations

import (
	"context"
	"fmt"
	"io/fs"
	"path"

	"github.com/thepabloaguilar/homelab/internal/ssh_operations/commands"
	"github.com/thepabloaguilar/homelab/internal/tools"
)

const sshdConfigPath = "/etc/ssh/sshd_config.d"

type hardenSSHSecurity struct {
	newConfigDir fs.FS
}

func NewHardenSSHSecurity(newConfigDir fs.FS) Operation {
	return &hardenSSHSecurity{
		newConfigDir: newConfigDir,
	}
}

func (s *hardenSSHSecurity) Do(ctx context.Context, client *tools.SSHClient) error {
	ctx = tools.NamedLoggerToContext(ctx, s.String())
	if err := s.cleanDefaultSSHConfigs(ctx, client); err != nil {
		return err
	}

	if err := s.pushNewSSHConfigs(ctx, client); err != nil {
		return err
	}

	if err := s.replaceDiffieHellmanModuli(ctx, client); err != nil {
		return err
	}

	if err := s.regenerateHostKeys(ctx, client); err != nil {
		return err
	}

	return nil
}

func (s *hardenSSHSecurity) cleanDefaultSSHConfigs(ctx context.Context, client *tools.SSHClient) error {
	_, err := client.Exec(ctx, fmt.Sprintf("/bin/rm %s", path.Join(sshdConfigPath, "*")))
	if err != nil {
		return err
	}

	return nil
}

func (s *hardenSSHSecurity) pushNewSSHConfigs(ctx context.Context, client *tools.SSHClient) error {
	return client.UploadDir(ctx, s.newConfigDir, ".")
}

func (s *hardenSSHSecurity) replaceDiffieHellmanModuli(ctx context.Context, client *tools.SSHClient) error {
	_, err := client.Exec(ctx, "awk '$5 >= 3071' /etc/ssh/moduli > /etc/ssh/moduli.safe")
	if err != nil {
		return err
	}

	_, err = client.Exec(ctx, "mv /etc/ssh/moduli.safe /etc/ssh/moduli")
	if err != nil {
		return err
	}

	return nil
}

func (s *hardenSSHSecurity) regenerateHostKeys(ctx context.Context, client *tools.SSHClient) error {
	_, err := client.Exec(ctx, "rm /etc/ssh/ssh_host_*")
	if err != nil {
		return err
	}

	keyGenCommands := []string{
		commands.SSHKeygen(
			"/etc/ssh/ssh_host_ed25519_key",
			commands.SSHKeyTypeED25519,
			0,
		),
		commands.SSHKeygen(
			"/etc/ssh/ssh_host_rsa_key",
			commands.SSHKeyTypeRSA,
			4096,
		),
	}

	for _, command := range keyGenCommands {
		_, err = client.Exec(ctx, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *hardenSSHSecurity) String() string {
	return "Clean SSH defaults"
}
