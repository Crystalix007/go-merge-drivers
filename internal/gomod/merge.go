package gomod

import (
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

// Merge merges the changes between the current and other go.mod files into the
// common ancestor go.mod file.
func Merge(current, other, ancestor modfile.File) modfile.File {
	currentChanges := Diff(current, ancestor)
	otherChanges := Diff(other, ancestor)

	mergedChanges := mergeChanges(currentChanges, otherChanges)

	// Now merge back into the ancestor file.
	return mergeChanges(mergedChanges, ancestor)
}

// mergeChanges merges the two changesets, preferring higher-versioned values,
// then the current changes over the other changes.
func mergeChanges(currentChanges, otherChanges modfile.File) modfile.File {
	otherReqs := make(map[string]modfile.Require)

	for _, req := range otherChanges.Require {
		otherReqs[req.Mod.Path] = *req
	}

	for _, req := range currentChanges.Require {
		otherReq, ok := otherReqs[req.Mod.Path]

		// If the other require statement doesn't exist, or the current
		// require is a higher version, then skip.
		if !ok || semver.Compare(req.Mod.Version, otherReq.Mod.Version) > 0 {
			otherChanges.AddRequire(req.Mod.Path, req.Mod.Version)
		}
	}

	for _, exc := range currentChanges.Exclude {
		otherChanges.AddExclude(exc.Mod.Path, exc.Mod.Version)
	}

	otherReps := make(map[string]modfile.Replace)

	for _, rep := range otherChanges.Replace {
		otherReps[rep.Old.Path] = *rep
	}

	for _, rep := range currentChanges.Replace {
		// Update the replace statement if the current replace is a higher
		// version than the existing one.
		if otherRep, ok := otherReps[rep.Old.Path]; !ok || semver.Compare(rep.New.Version, otherRep.New.Version) > 0 {
			otherChanges.AddReplace(rep.Old.Path, rep.Old.Version, rep.New.Path, rep.New.Version)
		}
	}

	return otherChanges
}
