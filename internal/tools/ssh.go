package tools

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"syscall"

	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"github.com/thepabloaguilar/homelab/internal/config"
)

type SSHClient struct {
	client *ssh.Client
}

func NewSSHClient(cfg config.ServerConfig) (*SSHClient, error) {
	authMethods := []ssh.AuthMethod{
		ssh.Password(cfg.Auth.Password),
	}

	if cfg.Auth.InteractivePassword {
		authMethods = append(authMethods, interactivePassword())
	}

	sshCfg := &ssh.ClientConfig{
		User: cfg.User,
		Auth: authMethods,
		// TODO!: We should not ignore an insecure host key
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         cfg.ConnectionTimeout,
	}

	client, err := ssh.Dial("tcp", cfg.Host, sshCfg)
	if err != nil {
		return nil, err
	}

	return &SSHClient{client: client}, nil
}

func (c *SSHClient) Exec(ctx context.Context, cmd string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer LogCloser(ctx, session)

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	if err != nil {
		LoggerFromContext(ctx).Error("error executing command",
			zap.String("cmd", cmd),
			zap.String("stderr", stderr.String()),
			zap.Error(err),
		)

		return "", err
	}

	LoggerFromContext(ctx).Info(
		"command executed",
		zap.String("cmd", cmd),
		zap.String("stdout", stdout.String()),
	)

	return stdout.String(), nil
}

func (c *SSHClient) UploadDir(ctx context.Context, dir fs.FS, root string) error {
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return err
	}
	defer LogCloser(ctx, sftpClient)

	err = fs.WalkDir(dir, root, func(from string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		to := filepath.Join("/", from)

		// We do not actually upload the whole dir in one shot, instead we upload one
		// file at time. We could use `rsync` but worthy right now.
		if d.IsDir() {
			return nil
		}

		LoggerFromContext(ctx).Info("copying file",
			zap.String("from", from),
			zap.String("to", to),
		)

		localFile, err := dir.Open(from)
		if err != nil {
			return err
		}
		defer LogCloser(ctx, localFile)

		remoteFile, err := sftpClient.Create(to)
		if err != nil {
			return err
		}
		defer LogCloser(ctx, remoteFile)

		_, err = io.Copy(remoteFile, localFile)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) Close() error {
	return c.client.Close()
}

func interactivePassword() ssh.AuthMethod {
	return ssh.PasswordCallback(func() (secret string, err error) {
		fmt.Printf("Enter password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}

		return string(password), nil
	})
}
