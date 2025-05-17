package gomod

import (
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

// Diff compares two modfile.File structs and returns the changes between them.
// It checks for differences in the require, exclude, and replace statements.
func Diff(version modfile.File, ancestor modfile.File) modfile.File {
	ancestorStr, err := ancestor.Format()
	if err != nil {
		panic(err)
	}

	changes, err := modfile.Parse(
		ancestor.Syntax.Name,
		[]byte(ancestorStr),
		nil,
	)
	if err != nil {
		panic(err)
	}

	semverVersionGoVersion := "v" + version.Go.Version
	semverAncestorGoVersion := "v" + ancestor.Go.Version

	semverComparison := semver.Compare(semverVersionGoVersion, semverAncestorGoVersion)

	// If the new version uses a later Go version, then update the Go version.
	if semverComparison > 0 {
		changes.Go = version.Go
	}

	// Clear the require, exclude, and replace statements.
	changes.Exclude = []*modfile.Exclude{}
	changes.Replace = []*modfile.Replace{}
	changes.Require = []*modfile.Require{}

	// Avoid quadratic behavior by creating a map of the ancestor require
	// statements.
	ancestorReqs := make(map[string]modfile.Require)

	for _, req := range ancestor.Require {
		ancestorReqs[req.Mod.Path] = *req
	}

	// Add the require statements.
	for _, req := range version.Require {
		if ancestorReq, ok := ancestorReqs[req.Mod.Path]; !ok || *req != ancestorReq {
			changes.Require = append(changes.Require, req)
		}
	}

	ancestorExcludes := make(map[string]modfile.Exclude)

	for _, exc := range ancestor.Exclude {
		ancestorExcludes[exc.Mod.Path] = *exc
	}

	// Add the exclude statements.
	for _, exc := range version.Exclude {
		if ancestorExc, ok := ancestorExcludes[exc.Mod.Path]; !ok || *exc != ancestorExc {
			changes.Exclude = append(changes.Exclude, exc)
		}
	}

	ancestorReps := make(map[string]modfile.Replace)

	for _, rep := range ancestor.Replace {
		ancestorReps[rep.Old.Path] = *rep
	}

	// Add the replace statements.
	for _, rep := range version.Replace {
		if ancestorRep, ok := ancestorReps[rep.Old.Path]; !ok || *rep != ancestorRep {
			changes.Replace = append(changes.Replace, rep)
		}
	}

	return *changes
}
