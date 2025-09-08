package commands

import (
	"strconv"
	"strings"
)

type SSHKeyType string

const (
	SSHKeyTypeED25519 SSHKeyType = "ed25519"
	SSHKeyTypeRSA     SSHKeyType = "rsa"
)

func SSHKeygen(
	path string,
	keyType SSHKeyType,
	keySize uint32,
) string {
	command := []string{
		"ssh-keygen",
		"-t", string(keyType),
		"-f", path,
		"-N", "\"\"",
	}

	if keySize > 0 {
		command = append(command, "-b", strconv.Itoa(int(keySize)))
	}

	return strings.Join(command, " ")
}
