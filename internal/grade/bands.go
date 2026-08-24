package grade

// This file contains the bookkeeping that keeps the strip table consistent:
// the bands must partition the positive axis with no gaps and no overlaps so
// that findBand always has exactly one answer.

// BandsPartition reports whether the pinned strip table covers every
// non-negative arc height exactly once. The first band must start at zero,
// every consecutive pair must meet edge to edge, and the last band must be
// open-ended.
func BandsPartition() bool {
	if len(Strips) == 0 {
		return false
	}
	if Strips[0].Lower != 0 {
		return false
	}
	for i := 1; i < len(Strips); i++ {
		if Strips[i].Lower != Strips[i-1].Upper {
			return false
		}
	}
	return Strips[len(Strips)-1].Upper == positiveInf()
}

// BandsOverlapFree reports whether no two bands share any arc height. It is
// the dual of the partition check and is asserted by the tests.
func BandsOverlapFree() bool {
	for i := 0; i < len(Strips); i++ {
		for j := i + 1; j < len(Strips); j++ {
			if Overlaps(Strips[i], Strips[j]) {
				return false
			}
		}
	}
	return true
}

// UniqueLetters reports whether every designator in the table appears exactly
// once. The designators are the public identifiers of the strips, so a
// duplicate would make recommendations ambiguous.
func UniqueLetters() bool {
	seen := make(map[string]bool)
	for _, s := range Strips {
		if seen[s.Letter] {
			return false
		}
		seen[s.Letter] = true
	}
	return true
}

// FirstBandEdge returns the lower edge of the first band and the upper edge of
// the last band, which is used by reports to describe the coverage of the
// table.
func FirstBandEdge() (float64, float64) {
	if len(Strips) == 0 {
		return 0, 0
	}
	return Strips[0].Lower, Strips[len(Strips)-1].Upper
}
