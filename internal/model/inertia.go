package model

import "math"

func MomentOfInertia(p Params) float64 {
	return p.Width * math.Pow(p.Thickness, 3) / 12.0
}

func MomentOfInertiaFor(width, thickness float64) float64 {
	return width * math.Pow(thickness, 3) / 12.0
}

func SectionModulus(p Params) float64 {
	return p.Width * p.Thickness * p.Thickness / 6.0
}

func NeutralAxisOffset(p Params) float64 {
	return p.Thickness / 2.0
}

func ThicknessCubed(thickness float64) float64 {
	return math.Pow(thickness, 3)
}
