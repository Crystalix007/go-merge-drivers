package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/crystalix007/go-merge-drivers/internal/flags"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "go-mod-merge",
	}

	flags := flags.AddFlags(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return run(cmd, flags)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		os.Exit(1)
	}
}

func run(cmd *cobra.Command, flags flags.Flags) error {
	slog.InfoContext(
		cmd.Context(),
		"running go.mod merge",
		slog.String("common-ancestor", *flags.CommonAncestor),
		slog.String("current-version", *flags.CurrentVersion),
		slog.String("other-version", *flags.OtherVersion),
	)

	return nil
}
