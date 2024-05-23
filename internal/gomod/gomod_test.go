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

// findReplaces returns a slice of modfile.Replace that match the given path.
func findReplaces(r modfile.File, oldPath string) []*modfile.Replace {
	var requires []*modfile.Replace

	for _, rep := range r.Replace {
		if rep.Old.Path == oldPath {
			requires = append(requires, rep)
		}
	}

	return requires
}
