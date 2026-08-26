package model

type StandardStrip struct {
	Letter    string
	Thickness float64
	Width     float64
	Length    float64
}

var StandardStrips = []StandardStrip{
	{Letter: "N", Thickness: 0.79, Width: 18.5, Length: 76.0},
	{Letter: "A", Thickness: 1.29, Width: 18.5, Length: 76.0},
	{Letter: "C", Thickness: 2.39, Width: 18.5, Length: 76.0},
}

func NominalFor(letter string) (StandardStrip, bool) {
	for _, s := range StandardStrips {
		if s.Letter == letter {
			return s, true
		}
	}
	return StandardStrip{}, false
}

func NominalThickness(letter string) (float64, bool) {
	s, ok := NominalFor(letter)
	if !ok {
		return 0, false
	}
	return s.Thickness, true
}

func (s StandardStrip) Matches(p Params) bool {
	const tol = 1e-6
	return close(s.Thickness, p.Thickness, tol) &&
		close(s.Width, p.Width, tol) &&
		close(s.Length, p.Length, tol)
}

func StripGeometryName(p Params) string {
	for _, s := range StandardStrips {
		if s.Matches(p) {
			return s.Letter + "-strip geometry"
		}
	}
	return "custom geometry"
}

func ThicknessDeviation(p Params, letter string) (float64, bool) {
	nominal, ok := NominalThickness(letter)
	if !ok {
		return 0, false
	}
	return p.Thickness - nominal, true
}

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
