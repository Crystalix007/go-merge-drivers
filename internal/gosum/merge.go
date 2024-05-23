package gosum

import "maps"

// Merge merges two go.sum files together. The current go.sum file is the one
// that is being modified, the other go.sum file is the one that is being merged
// in, and the ancestor go.sum file is the common ancestor of the two go.sum
// files.
//
// If there are inconsistent hashes between the current and other go.sum files,
// an error is returned.
func Merge(current GoSum, other GoSum, ancestor GoSum) (GoSum, error) {
	currentAdded, currentModified, currentRemoved := Diff(current, ancestor)
	otherAdded, otherModified, otherRemoved := Diff(other, ancestor)

	// Check for inconsistent hashes in the added / updated hashes.
	currentAddedModified := overlay(currentAdded, currentModified)
	otherAddedModified := overlay(otherAdded, otherModified)

	_, modified, _ := Diff(currentAddedModified, otherAddedModified)
	if len(modified) != 0 {
		return nil, ErrHashMismatch
	}

	allAddedModified := overlay(currentAddedModified, otherAddedModified)
	allRemoved := overlay(currentRemoved, otherRemoved)

	// Delete any new / modified hashes from the removal list, to keep
	// all potentially required hashes.
	for key := range allAddedModified {
		delete(allRemoved, key)
	}

	res := maps.Clone(ancestor)

	// Forcibly change the hashes in the ancestor go.sum file to the new hashes,
	// since we've agreed these updates are not conflicting with each other.
	for key, hash := range allAddedModified {
		res[key] = hash
	}

	// Remove any hashes that have been removed.
	for key := range allRemoved {
		delete(res, key)
	}

	return res, nil
}

// overlay overlays the sum go.sum file over the ancestor go.sum file.
func overlay(sum GoSum, ancestor GoSum) GoSum {
	overlay := maps.Clone(sum)

	for key, hash := range ancestor {
		if _, ok := overlay[key]; !ok {
			overlay[key] = hash
		}
	}

	return overlay
}
