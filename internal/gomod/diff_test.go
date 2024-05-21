package gomod_test

import (
	"slices"
	"testing"

	"github.com/crystalix007/go-merge-drivers/internal/gomod"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

func TestDiff_same(t *testing.T) {
	t.Parallel()

	modfile, err := gomod.Parse("testdata/current.go.mod")
	require.NoError(t, err)

	diff := gomod.Diff(*modfile, *modfile)

	require.Empty(t, diff.Require)
	require.Empty(t, diff.Exclude)
	require.Empty(t, diff.Replace)
}

func TestDiff_currentVsAncestor(t *testing.T) {
	t.Parallel()

	current, err := gomod.Parse("testdata/current.go.mod")
	require.NoError(t, err)

	ancestor, err := gomod.Parse("testdata/ancestor.go.mod")
	require.NoError(t, err)

	diff := gomod.Diff(*current, *ancestor)

	// New require statement.
	slices.ContainsFunc(diff.Require, func(r *modfile.Require) bool {
		return r.Mod.Path == "github.com/spf13/cobra"
	})
}

func TestDiff_currentVsOther(t *testing.T) {
	t.Parallel()

	current, err := gomod.Parse("testdata/current.go.mod")
	require.NoError(t, err)

	other, err := gomod.Parse("testdata/other.go.mod")
	require.NoError(t, err)

	diff := gomod.Diff(*current, *other)

	// Changed require statement.
	slices.ContainsFunc(diff.Require, func(r *modfile.Require) bool {
		return r.Mod.Path == "github.com/spf13/cobra"
	})

	// Removed replace statement.
	slices.ContainsFunc(diff.Replace, func(r *modfile.Replace) bool {
		return r.Old.Path == "github.com/spf13/cobra"
	})
}
