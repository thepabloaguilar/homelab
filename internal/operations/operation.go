package operations

import (
	"context"

	"github.com/thepabloaguilar/homelab/internal/config"
)

type Operation func(context.Context, config.ServerConfig) error
