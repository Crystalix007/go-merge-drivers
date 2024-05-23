package gosum_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/crystalix007/go-merge-drivers/internal/gosum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGoSum_empty(t *testing.T) {
	t.Parallel()

	emptyBytes := bytes.NewReader([]byte{})

	gosum, err := gosum.NewGoSum(emptyBytes)

	require.NoError(t, err)
	require.NotNil(t, gosum)
	assert.Empty(t, gosum)
}

func TestNewGoSum_singleMod(t *testing.T) {
	t.Parallel()

	modFile, err := os.Open("testdata/singlemod.go.sum")
	require.NoError(t, err)

	defer modFile.Close()

	goSum, err := gosum.NewGoSum(modFile)
	require.NoError(t, err)

	modKey := gosum.GoSumKey{
		ModulePath: "golang.org/x/mod",
		Version:    "v0.17.0",
		Path:       "",
	}

	require.Contains(t, goSum, modKey)
	assert.Equal(t, gosum.GoSumHash("h1:zY54UmvipHiNd+pm+m0x9KhZ9hl1/7QNMyxXbc6ICqA="), goSum[modKey])

	modKey.Path = "go.mod"

	require.Contains(t, goSum, modKey)
	assert.Equal(t, gosum.GoSumHash("h1:hTbmBsO62+eylJbnUtE2MGJUyE7QWk4xUqPFrRgJ+7c="), goSum[modKey])
}

func TestNewGoSum_duplicateMod(t *testing.T) {
	t.Parallel()

	modFile, err := os.Open("testdata/duplicatemod.go.sum")
	require.NoError(t, err)

	defer modFile.Close()

	_, err = gosum.NewGoSum(modFile)

	require.ErrorIs(t, err, gosum.ErrHashMismatch)
}

func TestGoSum_String(t *testing.T) {
	t.Parallel()

	goSum := gosum.GoSum{
		gosum.GoSumKey{
			ModulePath: "golang.org/x/mod",
			Version:    "v0.17.0",
			Path:       "",
		}: gosum.GoSumHash("h1:zY54UmvipHiNd+pm+m0x9KhZ9hl1/7QNMyxXbc6ICqA="),
		gosum.GoSumKey{
			ModulePath: "golang.org/x/mod",
			Version:    "v0.17.0",
			Path:       "go.mod",
		}: gosum.GoSumHash("h1:hTbmBsO62+eylJbnUtE2MGJUyE7QWk4xUqPFrRgJ+7c="),
		gosum.GoSumKey{
			ModulePath: "golang.org/x/exp",
			Version:    "v0.0.0-20240506185415-9bf2ced13842",
			Path:       "",
		}: "h1:vr/HnozRka3pE4EsMEg1lgkXJkTFJCVUX+S/ZT6wYzM=",
		gosum.GoSumKey{
			ModulePath: "golang.org/x/exp",
			Version:    "v0.0.0-20240506185415-9bf2ced13842",
			Path:       "go.mod",
		}: "h1:XtvwrStGgqGPLc4cjQfWqZHG1YFdYs6swckp8vpsjnc=",
	}

	expected := `
golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 h1:vr/HnozRka3pE4EsMEg1lgkXJkTFJCVUX+S/ZT6wYzM=
golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842/go.mod h1:XtvwrStGgqGPLc4cjQfWqZHG1YFdYs6swckp8vpsjnc=
golang.org/x/mod v0.17.0 h1:zY54UmvipHiNd+pm+m0x9KhZ9hl1/7QNMyxXbc6ICqA=
golang.org/x/mod v0.17.0/go.mod h1:hTbmBsO62+eylJbnUtE2MGJUyE7QWk4xUqPFrRgJ+7c=
	`
	actual := goSum.String()

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(actual))
}
