package grade

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

func FirstBandEdge() (float64, float64) {
	if len(Strips) == 0 {
		return 0, 0
	}
	return Strips[0].Lower, Strips[len(Strips)-1].Upper
}
