package gomod

import (
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

// Parse parses the contents of a go.mod file and returns a modfile.File
// object.
func Parse(filename string) (*modfile.File, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read go.mod file (%s): %w",
			filename,
			err,
		)
	}

	mod, err := modfile.Parse(filename, data, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse go.mod file (%s): %w",
			filename,
			err,
		)
	}

	mod.Cleanup()

	return mod, nil
}
