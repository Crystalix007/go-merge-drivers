package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/crystalix007/go-merge-drivers/internal/flags"
	"github.com/crystalix007/go-merge-drivers/internal/gomod"
	"github.com/crystalix007/go-merge-drivers/internal/gosum"
	"github.com/spf13/cobra"
)

var (
	// ErrUnknownFile is returned when the file is not a recognised go module
	// file.
	ErrUnknownFile = errors.New("file is not a go module file")

	// ErrNoCommonAncestor is returned when the common ancestor is not provided.
	ErrNoCommonAncestor = errors.New("common ancestor is not provided")

	// ErrNoCurrentVersion is returned when the current version is not provided.
	ErrNoCurrentVersion = errors.New("current version is not provided")

	// ErrNoOtherVersion is returned when the other version is not provided.
	ErrNoOtherVersion = errors.New("other version is not provided")

	// ErrNoResult is returned when the result is not provided.
	ErrNoResult = errors.New("result is not provided")
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

	var resultFile io.WriteCloser

	if *flags.Output == "/dev/stdout" {
		resultFile = os.Stdout
	} else {
		var err error

		resultFile, err = os.Create(*flags.Output)
		if err != nil {
			return fmt.Errorf(
				"failed to create output file (%s): %w",
				*flags.Output,
				err,
			)
		}
	}

	defer resultFile.Close()

	_, filename := path.Split(*flags.Result)

	switch filename {
	case "go.mod":
		return runGoModMerge(cmd.Context(), flags, resultFile)
	case "go.sum":
		return runGoSumMerge(flags, resultFile)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownFile, filename)
	}
}

// runGoModMerge will run the go.mod merge operation.
func runGoModMerge(ctx context.Context, flags flags.Flags, output io.Writer) error {
	slog.DebugContext(
		ctx,
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

	if _, err := output.Write(mergedBytes); err != nil {
		return fmt.Errorf(
			"failed to write go.mod file (%s): %w",
			*flags.Result,
			err,
		)
	}

	return nil
}

// runGoSumMerge will run the go.sum merge operation.
func runGoSumMerge(flags flags.Flags, output io.Writer) error {
	current, err := parseGoSumFile(*flags.CurrentVersion)
	if err != nil {
		return fmt.Errorf(
			"failed to parse current go.sum file: %w",
			err,
		)
	}

	other, err := parseGoSumFile(*flags.OtherVersion)
	if err != nil {
		return fmt.Errorf(
			"failed to parse other go.sum file: %w",
			err,
		)
	}

	ancestor, err := parseGoSumFile(*flags.CommonAncestor)
	if err != nil {
		return fmt.Errorf(
			"failed to parse common ancestor go.sum file: %w",
			err,
		)
	}

	merged, err := gosum.Merge(current, other, ancestor)
	if err != nil {
		return fmt.Errorf(
			"failed to merge go.sum files: %w",
			err,
		)
	}

	resultFile, err := os.Create(*flags.Result)
	if err != nil {
		return fmt.Errorf(
			"failed to create result go.sum file: %w",
			err,
		)
	}

	defer resultFile.Close()

	result := merged.String()

	if _, err := output.Write([]byte(result)); err != nil {
		return fmt.Errorf(
			"failed to write go.sum file (%s): %w",
			*flags.Result,
			err,
		)
	}

	return nil
}

func parseGoSumFile(path string) (gosum.GoSum, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open go.sum file (%s): %w",
			path,
			err,
		)
	}

	defer file.Close()

	goSum, err := gosum.NewGoSum(file)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse go.sum file: %w",
			err,
		)
	}

	return goSum, nil
}

// check will verify the provided flags.
func check(flags flags.Flags) error {
	if *flags.CommonAncestor == "" {
		return ErrNoCommonAncestor
	}

	if *flags.CurrentVersion == "" {
		return ErrNoCurrentVersion
	}

	if *flags.OtherVersion == "" {
		return ErrNoOtherVersion
	}

	if *flags.Result == "" {
		return ErrNoResult
	}

	return nil
}
