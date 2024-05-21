package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/crystalix007/go-merge-drivers/internal/flags"
	"github.com/crystalix007/go-merge-drivers/internal/gomod"
	"github.com/spf13/cobra"
)

var (
	// ErrNotGoMod is returned when the file is not a go.mod file.
	ErrNotGoMod = errors.New("file is not a go.mod file")

	// ErrNoCommonAncestor is returned when the common ancestor is not provided.
	ErrNoCommonAncestor = errors.New("common ancestor is not provided")

	// ErrNoCurrentVersion is returned when the current version is not provided.
	ErrNoCurrentVersion = errors.New("current version is not provided")

	// ErrNoOtherVersion is returned when the other version is not provided.
	ErrNoOtherVersion = errors.New("other version is not provided")
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
	if err := check(flags); err != nil {
		return err
	}

	slog.InfoContext(
		cmd.Context(),
		"running go.mod merge",
		slog.String("common-ancestor", *flags.CommonAncestor),
		slog.String("current-version", *flags.CurrentVersion),
		slog.String("other-version", *flags.OtherVersion),
		slog.String("result", *flags.Result),
	)

	commonAncestor, err := gomod.Parse(*flags.CommonAncestor)
	if err != nil {
		return fmt.Errorf(
			"failed to parse common ancestor: %w",
			err,
		)
	}

	currentVersion, err := gomod.Parse(*flags.CurrentVersion)
	if err != nil {
		return fmt.Errorf(
			"failed to parse current version: %w",
			err,
		)
	}

	otherVersion, err := gomod.Parse(*flags.OtherVersion)
	if err != nil {
		return fmt.Errorf(
			"failed to parse other version: %w",
			err,
		)
	}

	// Merge the go.mod file changes.
	merged := gomod.Merge(*currentVersion, *otherVersion, *commonAncestor)

	merged.Cleanup()

	mergedBytes, err := merged.Format()
	if err != nil {
		return fmt.Errorf(
			"failed to format go.mod file: %w",
			err,
		)
	}

	if err := os.WriteFile(*flags.Result, mergedBytes, 0644); err != nil {
		return fmt.Errorf(
			"failed to write go.mod file (%s): %w",
			*flags.Result,
			err,
		)
	}

	return nil
}

// check will verify the provided flags.
func check(flags flags.Flags) error {
	_, filename := path.Split(*flags.Result)
	if filename != "go.mod" {
		return ErrNotGoMod
	}

	if *flags.CommonAncestor == "" {
		return ErrNoCommonAncestor
	}

	if *flags.CurrentVersion == "" {
		return ErrNoCurrentVersion
	}

	if *flags.OtherVersion == "" {
		return ErrNoOtherVersion
	}

	return nil
}
