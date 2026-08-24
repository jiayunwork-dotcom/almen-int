package grade

// Strip describes one Almen strip designator and the band of arc heights it
// is recommended for. All lengths are in millimetres.
type Strip struct {
	Letter    string  // designator: "N", "A" or "C"
	Thickness float64 // nominal strip thickness, mm
	Lower     float64 // lower edge of the recommended band, inclusive, mm
	Upper     float64 // upper edge of the recommended band, exclusive, mm
	Note      string  // one-line description of the strip's role
}

// Strips is the pinned Almen strip table. The bands are non-overlapping and
// partition the positive arc-height axis: N covers the low end, A covers the
// standard medium band, and C covers the high end.
var Strips = []Strip{
	{
		Letter:    "N",
		Thickness: 0.79,
		Lower:     0.00,
		Upper:     0.10,
		Note:      "thin strip, low intensity",
	},
	{
		Letter:    "A",
		Thickness: 1.29,
		Lower:     0.10,
		Upper:     0.60,
		Note:      "standard strip, medium intensity",
	},
	{
		Letter:    "C",
		Thickness: 2.39,
		Lower:     0.60,
		Upper:     positiveInf(),
		Note:      "thick strip, high intensity",
	},
}

// ByLetter returns the strip with the given designator, or false when the
// letter is unknown.
func ByLetter(letter string) (Strip, bool) {
	for _, s := range Strips {
		if s.Letter == letter {
			return s, true
		}
	}
	return Strip{}, false
}

// BandContains reports whether the arc height falls inside the strip's
// recommended band. The lower edge is inclusive and the upper edge exclusive,
// which makes the bands partition the axis without gaps or overlaps.
func (s Strip) BandContains(h float64) bool {
	return h >= s.Lower && h < s.Upper
}

// BandWidth returns the width of the recommended band in millimetres. The C
// strip is open-ended, so its width is reported as positive infinity.
func (s Strip) BandWidth() float64 {
	return s.Upper - s.Lower
}

// BandCenter returns the midpoint of a finite band, or the lower edge when the
// band is open-ended.
func (s Strip) BandCenter() float64 {
	if s.Upper == positiveInf() {
		return s.Lower
	}
	return (s.Lower + s.Upper) / 2
}

// Overlaps reports whether two bands share any arc height. The pinned table is
// expected to have no overlaps; the check exists so tests can assert it.
func Overlaps(a, b Strip) bool {
	return a.Lower < b.Upper && b.Lower < a.Upper
}

func positiveInf() float64 {
	return 1e300 // a large sentinel that behaves as infinity for band math
}
