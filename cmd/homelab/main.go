package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
)

func main() {
	ctx := context.Background()
	if err := fang.Execute(ctx, rootCmd); err != nil {
		os.Exit(1)
	}
}
