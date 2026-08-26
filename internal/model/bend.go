package model

import "math"

func ModulusPerArea(p Params) float64 {
	return p.Modulus * 1000.0
}

func Curvature(p Params) float64 {
	return BendingMoment(p) / (ModulusPerArea(p) * MomentOfInertia(p))
}

func ArcRadius(p Params) float64 {
	return 1.0 / Curvature(p)
}

func PlateauArcHeight(p Params) float64 {
	kappa := Curvature(p)
	return kappa * p.Length * p.Length / 8.0
}

func ExactSagitta(p Params) float64 {
	R := ArcRadius(p)
	half := p.Length / 2.0
	return R - math.Sqrt(R*R-half*half)
}

func RelativeError(p Params) float64 {
	h := PlateauArcHeight(p)
	exact := ExactSagitta(p)
	if exact == 0 {
		return 0
	}
	return math.Abs(exact-h) / exact
}
