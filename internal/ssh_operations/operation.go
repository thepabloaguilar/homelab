package ssh_operations

import (
	"context"
	"fmt"

	"github.com/thepabloaguilar/homelab/internal/tools"
)

type Operation interface {
	fmt.Stringer
	Do(ctx context.Context, client *tools.SSHClient) error
}
