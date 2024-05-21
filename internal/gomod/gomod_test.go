package gomod_test

import "golang.org/x/mod/modfile"

// findRequires returns a slice of modfile.Require that match the given path.
func findRequires(r modfile.File, path string) []*modfile.Require {
	var requires []*modfile.Require

	for _, req := range r.Require {
		if req.Mod.Path == path {
			requires = append(requires, req)
		}
	}

	return requires
}
