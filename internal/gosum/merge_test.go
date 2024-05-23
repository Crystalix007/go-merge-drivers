package gosum_test

import (
	"os"
	"testing"

	"github.com/crystalix007/go-merge-drivers/internal/gosum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	current := loadGoSum(t, "testdata/current.go.sum")
	other := loadGoSum(t, "testdata/other.go.sum")
	ancestor := loadGoSum(t, "testdata/ancestor.go.sum")
	expected := loadGoSum(t, "testdata/merged.go.sum")

	merged, err := gosum.Merge(current, other, ancestor)
	require.NoError(t, err)

	assert.Equal(t, expected, merged)
}

func loadGoSum(t *testing.T, path string) gosum.GoSum {
	t.Helper()

	file, err := os.Open(path)
	require.NoError(t, err)

	defer file.Close()

	goSum, err := gosum.NewGoSum(file)
	require.NoError(t, err)

	return goSum
}
