package model

func BendingMoment(p Params) float64 {
	return ResidualStress(p) * p.Width * p.LayerDepth * (p.Thickness - p.LayerDepth) / 2.0
}

func BendingMomentAtStress(p Params, sigma float64) float64 {
	return sigma * p.Width * p.LayerDepth * (p.Thickness - p.LayerDepth) / 2.0
}

func MomentPerUnitWidth(p Params) float64 {
	return BendingMoment(p) / p.Width
}

func LeverArm(p Params) float64 {
	return (p.Thickness - p.LayerDepth) / 2.0
}

func LayerForce(p Params) float64 {
	return ResidualStress(p) * p.Width * p.LayerDepth
}
