package model

import "math"

func ResidualStress(p Params) float64 {
	r := VelocityRatio(p)
	return p.ResidualStress * math.Pow(r, velocityPower)
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
