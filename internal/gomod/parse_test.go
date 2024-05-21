package gomod_test

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"

	"github.com/crystalix007/go-merge-drivers/internal/gomod"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

func TestParse_missingFile(t *testing.T) {
	t.Parallel()

	_, err := gomod.Parse("missing_file_" + randHex(t))
	require.Error(t, err)
}

func TestParse_empty(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp("", "go.mod")
	require.NoError(t, err)

	defer f.Close()
	defer os.Remove(f.Name())

	parsed, err := gomod.Parse(f.Name())
	require.NoError(t, err)
	require.NotNil(t, parsed)
}

func TestParse_current(t *testing.T) {
	t.Parallel()

	parsed, err := gomod.Parse("testdata/current.go.mod")
	require.NoError(t, err)

	require.NotNil(t, parsed)

	var cobraRequire *modfile.Require

	for _, req := range parsed.Require {
		if req.Mod.Path == "github.com/spf13/cobra" {
			cobraRequire = req
			break
		}
	}

	require.NotNil(t, cobraRequire)
	assert.Equal(t, "v1.8.0", cobraRequire.Mod.Version)
	assert.False(t, cobraRequire.Indirect)
}

func randHex(t *testing.T) string {
	t.Helper()

	var bytes [12]byte

	_, err := rand.Read(bytes[:])
	require.NoError(t, err)

	return hex.EncodeToString(bytes[:])
}
