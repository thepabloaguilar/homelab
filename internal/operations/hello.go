package operations

import (
	"context"

	"go.uber.org/zap"

	"github.com/thepabloaguilar/homelab/internal/config"
	"github.com/thepabloaguilar/homelab/internal/tools"
)

func Hello() Operation {
	return func(ctx context.Context, sc config.ServerConfig) error {
		client, err := tools.NewClient(sc)
		if err != nil {
			return err
		}

		output, err := client.Exec(ctx, "echo 'hello world'")
		if err != nil {
			return err
		}

		tools.LoggerFromContext(ctx).Info("command executed successfully",
			zap.String("output", output),
		)

		return nil
	}
}
