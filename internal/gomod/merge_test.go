package gomod_test

import (
	"testing"

	"github.com/crystalix007/go-merge-drivers/internal/gomod"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	current, err := gomod.Parse("testdata/current.go.mod")
	require.NoError(t, err)

	other, err := gomod.Parse("testdata/other.go.mod")
	require.NoError(t, err)

	ancestor, err := gomod.Parse("testdata/ancestor.go.mod")
	require.NoError(t, err)

	merged := gomod.Merge(*current, *other, *ancestor)

	cobraRequires := findRequires(merged, "github.com/spf13/cobra")
	require.Len(t, cobraRequires, 1)

	// Use current.go.mod version, as it is greater.
	assert.Equal(t, "v1.8.0", cobraRequires[0].Mod.Version)

	mousetrapRequires := findRequires(merged, "github.com/inconshreveable/mousetrap")
	require.Len(t, mousetrapRequires, 1)

	// Use other.go.mod version, as it is greater.
	assert.Equal(t, "v1.2.0", mousetrapRequires[0].Mod.Version)

	pflagRequires := findRequires(merged, "github.com/spf13/pflag")
	require.Len(t, pflagRequires, 1)

	// Use current.go.mod version, as it is greater.
	assert.Equal(t, "v1.0.5", pflagRequires[0].Mod.Version)
	assert.False(t, pflagRequires[0].Indirect)

	modRequires := findRequires(merged, "golang.org/x/mod")
	require.Len(t, modRequires, 1)

	// Use version from other.go.mod, as it is a full major version greater.
	assert.Equal(t, "v1.17.0", modRequires[0].Mod.Version)
}
