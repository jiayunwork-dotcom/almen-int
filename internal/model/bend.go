package model

import "math"

// ModulusPerArea converts the elastic modulus from gigapascals to N/mm^2.
// The bending formula needs stress-like units, and the case files carry the
// modulus in GPa because that is how material data sheets quote it.
func ModulusPerArea(p Params) float64 {
	return p.Modulus * 1000.0
}

// Curvature returns the deformed curvature of the strip, in 1/mm, from the
// thin-beam relation:
//
//	kappa = M / (E * I)
//
// The curvature is the linear bridge between the applied moment and the strip
// stiffness. All three factors (moment, modulus, inertia) are independent
// inputs, which is what lets the cross rules isolate each one.
func Curvature(p Params) float64 {
	return BendingMoment(p) / (ModulusPerArea(p) * MomentOfInertia(p))
}

// ArcRadius returns the radius of curvature 1/kappa in millimetres. It is the
// geometric counterpart of the curvature and is useful when reasoning about
// whether the small-deflection sagitta formula stays accurate.
func ArcRadius(p Params) float64 {
	return 1.0 / Curvature(p)
}

// PlateauArcHeight returns the arc height of a fully peened strip, in mm, for
// the small-deflection sagitta of a circular arc over a span L:
//
//	h = kappa * L^2 / 8
//
// This is the pinned bending result h proportional to M*L^2/(E*I) from the
// project contract. The gain applied for partial coverage is applied by the
// coverage package, never here.
func PlateauArcHeight(p Params) float64 {
	kappa := Curvature(p)
	return kappa * p.Length * p.Length / 8.0
}

// ExactSagitta returns the exact circular-arc sagitta
//
//	h = R - sqrt(R^2 - (L/2)^2)
//
// for comparison with the small-deflection value. For realistic Almen strips
// the deflection is a small fraction of the span, so the two agree to many
// significant digits; the exact value is kept as a reference, not as the
// reported arc height.
func ExactSagitta(p Params) float64 {
	R := ArcRadius(p)
	half := p.Length / 2.0
	return R - math.Sqrt(R*R-half*half)
}

// RelativeError returns the relative difference between the pinned arc height
// and the exact sagitta. It is used by reports to state how close the
// small-deflection formula is for a given case.
func RelativeError(p Params) float64 {
	h := PlateauArcHeight(p)
	exact := ExactSagitta(p)
	if exact == 0 {
		return 0
	}
	return math.Abs(exact-h) / exact
}
