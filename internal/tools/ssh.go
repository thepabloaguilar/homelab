package tools

import (
	"bytes"
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/thepabloaguilar/homelab/internal/config"
)

type Client struct {
	client *ssh.Client
}

func NewClient(cfg config.ServerConfig) (*Client, error) {
	sshCfg := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Auth.Password),
		},
		// TODO!: We should not ignore an insecure host key
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         cfg.ConnectionTimeout,
	}

	client, err := ssh.Dial("tcp", cfg.Host, sshCfg)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (s *Client) Exec(ctx context.Context, cmd string) (string, error) {
	session, err := s.client.NewSession()
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

	LoggerFromContext(ctx).Debug(
		"command executed",
		zap.String("cmd", cmd),
		zap.String("stdout", stdout.String()),
	)

	return stdout.String(), nil
}
