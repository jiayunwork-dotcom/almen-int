package model

import "math"

// MomentOfInertia returns the second moment of area of the rectangular strip
// cross-section about the neutral axis parallel to the width, in mm^4:
//
//	I = w * t^3 / 12
//
// The thickness enters at the third power. This cubic dependence is pinned:
// it is the whole reason a thicker strip bends less for the same moment, and
// a model that made the inertia linear in the thickness would remove that
// behaviour. The cross tests assert the cube directly.
func MomentOfInertia(p Params) float64 {
	return p.Width * math.Pow(p.Thickness, 3) / 12.0
}

// MomentOfInertiaFor evaluates the rectangular inertia for an explicit width
// and thickness. It is used by the sensitivity helpers to compare two strip
// sizes without building two full parameter sets.
func MomentOfInertiaFor(width, thickness float64) float64 {
	return width * math.Pow(thickness, 3) / 12.0
}

// SectionModulus returns the elastic section modulus I / (t/2) = w*t^2/6 in
// mm^3. It is the ratio of the bending moment to the peak bending stress at
// the strip surfaces and is reported as an intermediate quantity.
func SectionModulus(p Params) float64 {
	return p.Width * p.Thickness * p.Thickness / 6.0
}

// NeutralAxisOffset returns half the strip thickness, the distance from the
// neutral axis to either surface, in millimetres.
func NeutralAxisOffset(p Params) float64 {
	return p.Thickness / 2.0
}

// ThicknessCubed returns t^3 for an explicit thickness. Keeping the cube in
// one place makes it impossible for the reports and the cross tests to
// disagree about how thickness scales the inertia.
func ThicknessCubed(thickness float64) float64 {
	return math.Pow(thickness, 3)
}
