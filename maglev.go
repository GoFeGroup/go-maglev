package maglev

import (
	"runtime"
	"sort"
)

// GetLookupTable returns the Maglev lookup table of the size "m" for the given
// backends. The lookup table contains the IDs of the given backends.
//
// Maglev algorithm might produce different lookup table for the same
// set of backends listed in a different order. To avoid that sort
// backends by name, as the names are the same on all nodes (in opposite
// to backend IDs which are node-local).
//
// The weights implementation is inspired by https://github.com/envoyproxy/envoy/pull/2982.
//
// A backend weight is honored by altering the frequency how often a backend's turn is
// selected.
// A backend weight is multiplied in each turn by (n + 1) and compared to
// weightCntr[backendName] value which is an incrementation of weightSum (but starts at
// backend's weight / number of backends, so that each backend is selected at least once). If this is lower
// than weightCntr[backendName], another backend has a turn (and weightCntr[backendName]
// is incremented). This way we honor the weights.
func GetLookupTable(backendsMap map[string]*Backend, m uint64) []int {
	if len(backendsMap) == 0 {
		return nil
	}

	backends := make([]string, 0, len(backendsMap))
	weightCntr := make(map[string]float64, len(backendsMap))
	weightSum := uint64(0)

	l := len(backendsMap)

	for name, b := range backendsMap {
		backends = append(backends, name)
		weightSum += uint64(b.Weight)
		weightCntr[name] = float64(b.Weight) / float64(l)
	}

	sort.Strings(backends)

	perm := getPermutation(backends, m, runtime.NumCPU())
	next := make([]int, len(backends))
	entry := make([]int, m)

	for j := uint64(0); j < m; j++ {
		entry[j] = -1
	}

	for n := uint64(0); n < m; n++ {
		i := int(n) % l
		for {
			// change the default selection of backend turns only if weights are used
			if weightSum/uint64(l) > 1 {
				if ((n + 1) * uint64(backendsMap[backends[i]].Weight)) <= uint64(weightCntr[backends[i]]) {
					i = (i + 1) % l
					continue
				}
				weightCntr[backends[i]] += float64(weightSum)
			}
			c := perm[i*int(m)+next[i]]
			for entry[c] >= 0 {
				next[i] += 1
				c = perm[i*int(m)+next[i]]
			}
			entry[c] = int(backendsMap[backends[i]].ID)
			next[i] += 1
			break
		}
	}
	return entry
}
