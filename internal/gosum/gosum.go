package gosum

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/mod/semver"
)

var (
	// ErrHashMismatch is returned when a hash is added to a go.sum file that does
	// not match the existing hash.
	ErrHashMismatch = errors.New("gosum: hash mismatch")

	// ErrGoSumMustHaveThreeFields is returned when a go.sum file invalidly has
	// more than 3 fields per line.
	ErrGoSumMustHaveThreeFields = errors.New(
		"gosum: go.sum must have 3 fields per line",
	)
)

// GoSumKey represents a key in a go.sum file. All sums should be uniquely identified by a GoSumKey.
type GoSumKey struct {
	ModulePath string
	Version    string
	Path       string
}

// CompareKeys compares two GoSumKeys.
//
// Performs lexicographical ordering of module paths, then semver comparison of
// versions, and finally lexicographical ordering of paths.
func CompareKeys(this, other GoSumKey) int {
	return cmp.Or(
		cmp.Compare(this.ModulePath, other.ModulePath),
		semver.Compare(this.Version, other.Version),
		cmp.Compare(this.Path, other.Path),
	)
}

// GoSumHash represents a hash in a go.sum file.
type GoSumHash string

// GoSum represents a go.sum file.
type GoSum map[GoSumKey]GoSumHash

// NewGoSum creates a new GoSum from the given byte slice.
func NewGoSum(r io.Reader) (GoSum, error) {
	scanner := bufio.NewScanner(r)
	sum := make(GoSum)

	for scanner.Scan() {
		line := scanner.Text()
		flags := strings.Fields(line)

		if len(flags) != 3 {
			return nil, ErrGoSumMustHaveThreeFields
		}

		version, path, _ := strings.Cut(flags[1], "/")

		key := GoSumKey{
			ModulePath: flags[0],
			Version:    version,
			Path:       path,
		}

		if err := sum.Add(key, GoSumHash(flags[2])); err != nil {
			return nil, err
		}
	}

	return sum, nil
}

// Add implements a hash-safe way to add a key-value pair to a GoSum map.
func (g GoSum) Add(key GoSumKey, hash GoSumHash) error {
	if existingHash, ok := g[key]; ok && existingHash != hash {
		return ErrHashMismatch
	}

	g[key] = hash

	return nil
}

// Ensure [GoSum] implements the [fmt.Stringer] interface.
var _ fmt.Stringer = GoSum{}

// String returns a string representation of the GoSum map.
// The format of the string is "<module path> <version> <hash>\n" for each
// key-value pair in the map.
// If the key has a non-empty path, it is appended to the version with a "/"
// separator.
// The key-value pairs are concatenated in alphabetical order.
func (g GoSum) String() string {
	var b strings.Builder

	keys := maps.Keys(g)
	slices.SortFunc(keys, CompareKeys)

	for _, key := range keys {
		version := key.Version
		hash := g[key]

		if key.Path != "" {
			version += "/" + key.Path
		}

		fmt.Fprintf(&b, "%s %s %s\n", key.ModulePath, version, hash)
	}

	return b.String()
}
