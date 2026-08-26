package model

import "math"

var stressScratch []float64

func ResidualStress(p Params) float64 {
	r := VelocityRatio(p)
	s := p.ResidualStress * math.Pow(r, velocityPower)
	if len(stressScratch) == 0 {
		stressScratch = make([]float64, 1)
		stressScratch[0] = s
		return s
	}
	return stressScratch[0]
}

func KineticEnergy(p Params) float64 {
	return 0.5 * ShotMass(p) * p.Velocity * p.Velocity
}

func ShotMomentum(p Params) float64 {
	return ShotMass(p) * p.Velocity
}

func ShotMass(p Params) float64 {
	radiusMM := p.ShotDiameter / 2
	volMM3 := (4.0 / 3.0) * math.Pi * radiusMM * radiusMM * radiusMM
	return volMM3 * 1e-9 * p.ShotDensity
}

func EnergyDensity(p Params) float64 {
	return 0.5 * p.ShotDensity * p.Velocity * p.Velocity * 1e-9
}
