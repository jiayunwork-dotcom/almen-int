package model

const velocityPower = 2.0

type Params struct {
	Velocity          float64
	ReferenceVelocity float64
	ShotDiameter      float64
	ShotDensity       float64
	Thickness         float64
	Width             float64
	Length            float64
	Modulus           float64
	ResidualStress    float64
	LayerDepth        float64
}

func VelocityRatio(p Params) float64 {
	return p.Velocity / p.ReferenceVelocity
}

func VelocityExponent() float64 {
	return velocityPower
}

func ThicknessOf(p Params) float64 {
	return p.Thickness
}

func ModulusOf(p Params) float64 {
	return p.Modulus
}

func LayerDepthOf(p Params) float64 {
	return p.LayerDepth
}
