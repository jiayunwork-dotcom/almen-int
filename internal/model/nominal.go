package model

// StandardStrip describes the fixed geometry of a calibrated Almen strip.
// Every Almen strip shares the same span and width; the three designators
// differ only in their thickness, which is what selects the band of intensity
// they measure.
type StandardStrip struct {
	Letter    string  // designator: "N", "A" or "C"
	Thickness float64 // mm
	Width     float64 // mm
	Length    float64 // mm
}

// StandardStrips is the pinned table of calibrated strip geometries. The span
// and the width are the same for all three; only the thickness changes. A
// case that uses these dimensions is called a nominal-strip case.
var StandardStrips = []StandardStrip{
	{Letter: "N", Thickness: 0.79, Width: 18.5, Length: 76.0},
	{Letter: "A", Thickness: 1.29, Width: 18.5, Length: 76.0},
	{Letter: "C", Thickness: 2.39, Width: 18.5, Length: 76.0},
}

// NominalFor returns the standard strip geometry for a designator, or false
// when the letter is unknown.
func NominalFor(letter string) (StandardStrip, bool) {
	for _, s := range StandardStrips {
		if s.Letter == letter {
			return s, true
		}
	}
	return StandardStrip{}, false
}

// NominalThickness returns the standard thickness for a designator. It is a
// convenience used by the geometry checks and the reports.
func NominalThickness(letter string) (float64, bool) {
	s, ok := NominalFor(letter)
	if !ok {
		return 0, false
	}
	return s.Thickness, true
}

// Matches reports whether a parameter set uses the exact geometry of this
// standard strip within a small tolerance. The tolerance absorbs floating
// point noise from the JSON decoding, not real geometry differences.
func (s StandardStrip) Matches(p Params) bool {
	const tol = 1e-6
	return close(s.Thickness, p.Thickness, tol) &&
		close(s.Width, p.Width, tol) &&
		close(s.Length, p.Length, tol)
}

// StripGeometryName returns a short label for the strip geometry of a
// parameter set: the designator when it matches a standard strip exactly,
// otherwise "custom".
func StripGeometryName(p Params) string {
	for _, s := range StandardStrips {
		if s.Matches(p) {
			return s.Letter + "-strip geometry"
		}
	}
	return "custom geometry"
}

// ThicknessDeviation returns the signed difference between the case thickness
// and the nominal thickness of the given designator, in millimetres. It lets a
// report tell the user how far a strip is from the calibrated size.
func ThicknessDeviation(p Params, letter string) (float64, bool) {
	nominal, ok := NominalThickness(letter)
	if !ok {
		return 0, false
	}
	return p.Thickness - nominal, true
}

// StripThicknessRatio returns the ratio of the case thickness to the nominal
// thickness of the given designator. A ratio of one means a calibrated strip.
func StripThicknessRatio(p Params, letter string) (float64, bool) {
	nominal, ok := NominalThickness(letter)
	if !ok || nominal == 0 {
		return 0, false
	}
	return p.Thickness / nominal, true
}

func close(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}
