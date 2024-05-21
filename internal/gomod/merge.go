package gomod

import "golang.org/x/mod/modfile"

// Merge merges the changes between the current and other go.mod files into the
// common ancestor go.mod file.
func Merge(current, other, ancestor modfile.File) modfile.File {
	currentChanges := Diff(current, ancestor)
	otherChanges := Diff(other, ancestor)

	mergedChanges := mergeChanges(currentChanges, otherChanges)

	// Now merge back into the ancestor file.
	return mergeChanges(mergedChanges, ancestor)
}

// mergeChanges merges the two changesets, preferring the current changes over
// the other changes.
func mergeChanges(currentChanges, otherChanges modfile.File) modfile.File {
	for _, req := range currentChanges.Require {
		otherChanges.AddRequire(req.Mod.Path, req.Mod.Version)
	}

	for _, exc := range currentChanges.Exclude {
		otherChanges.AddExclude(exc.Mod.Path, exc.Mod.Version)
	}

	for _, rep := range currentChanges.Replace {
		otherChanges.AddReplace(rep.Old.Path, rep.Old.Version, rep.New.Path, rep.New.Version)
	}

	return otherChanges
}
