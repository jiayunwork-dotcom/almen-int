package model

import (
	"errors"
	"math"
)

func TriangleLayerForce(p Params) float64 {
	return ResidualStress(p) * p.Width * p.LayerDepth / 2.0
}

func TriangleCentroidFromSurface(p Params) float64 {
	return p.LayerDepth / 3.0
}

func TriangleLeverArm(p Params) float64 {
	return NeutralAxisOffset(p) - TriangleCentroidFromSurface(p)
}

func TriangleBendingMoment(p Params) float64 {
	arm := TriangleLeverArm(p)
	if arm <= 0 {
		return 0
	}
	return TriangleLayerForce(p) * arm
}

func TriangleCurvature(p Params) float64 {
	den := ModulusPerArea(p) * MomentOfInertia(p)
	if den == 0 {
		return 0
	}
	return TriangleBendingMoment(p) / den
}

func TrianglePlateauArcHeight(p Params) float64 {
	return TriangleCurvature(p) * p.Length * p.Length / 8.0
}

func ProfileMomentRatio(p Params) float64 {
	u := BendingMoment(p)
	if u == 0 {
		return 0
	}
	return TriangleBendingMoment(p) / u
}

func UniformEquivalentPeak(p Params) (float64, error) {
	u := BendingMoment(p)
	unit := TriangleBendingMoment(p)
	peak := ResidualStress(p)
	if peak == 0 || unit == 0 {
		return 0, errors.New("model: cannot scale a zero triangle moment")
	}
	return peak * (u / unit), nil
}

func TriangleLeverPositive(p Params) bool {
	return TriangleLeverArm(p) > 0
}

func ProfileForceRatio(p Params) float64 {
	u := LayerForce(p)
	if u == 0 {
		return 0
	}
	return TriangleLayerForce(p) / u
}

func Abs(x float64) float64 {
	return math.Abs(x)
}
