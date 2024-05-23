package gosum

// Diff returns the difference between the sum and ancestor go.sum files.
//
// Returns the added, modified, and removed hashes.
func Diff(sum GoSum, ancestor GoSum) (GoSum, GoSum, GoSum) {
	added := make(GoSum)
	modified := make(GoSum)
	removed := make(GoSum)

	for key, hash := range sum {
		if ancestorHash, ok := ancestor[key]; !ok {
			added[key] = hash
		} else if hash != ancestorHash {
			modified[key] = hash
		}
	}

	for key, hash := range ancestor {
		if _, ok := sum[key]; !ok {
			removed[key] = hash
		}
	}

	return added, modified, removed
}
