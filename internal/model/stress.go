package model

import "math"

// ResidualStress returns the compressive stress induced by the shot stream at
// the case velocity. The pinned power law scales the stress quoted at the
// reference velocity by the square of the velocity ratio:
//
//	sigma(v) = sigma_ref * (v / v_ref)^2
//
// The square mirrors the kinetic energy of a shot, which is proportional to
// v^2 for a fixed shot mass. Doubling the velocity therefore quadruples both
// the delivered energy density and the residual stress.
func ResidualStress(p Params) float64 {
	r := VelocityRatio(p)
	return p.ResidualStress * math.Pow(r, velocityPower)
}

// KineticEnergy returns the kinetic energy of a single shot in joules:
//
//	KE = (1/2) * m * v^2
//
// It is reported as a cross-check of the pinned velocity power law and is used
// nowhere inside the bending moment itself; a bug that forgets the square
// factor is caught by the kinetic-energy invariant tests.
func KineticEnergy(p Params) float64 {
	return 0.5 * ShotMass(p) * p.Velocity * p.Velocity
}

// ShotMomentum returns the linear momentum m*v of a single shot in kg*m/s.
// Momentum is linear in velocity and therefore grows more slowly than the
// kinetic energy; exposing both lets a report show why the model pins the
// energy rather than the momentum.
func ShotMomentum(p Params) float64 {
	return ShotMass(p) * p.Velocity
}

// ShotMass returns the mass of a single spherical shot in kilograms. The
// diameter is given in millimetres, so the volume is converted from cubic
// millimetres to cubic metres with the factor 1e-9 before it is multiplied by
// the density.
func ShotMass(p Params) float64 {
	radiusMM := p.ShotDiameter / 2
	volMM3 := (4.0 / 3.0) * math.Pi * radiusMM * radiusMM * radiusMM
	return volMM3 * 1e-9 * p.ShotDensity
}

// EnergyDensity returns the kinetic energy per unit shot volume in
// J/mm^3. Because both the mass and the energy carry the same shot volume,
// the density comparison cancels the volume and leaves v^2 * rho / 2; the
// value is provided so reports can compare two peening settings without
// recomputing the shot geometry.
func EnergyDensity(p Params) float64 {
	return 0.5 * p.ShotDensity * p.Velocity * p.Velocity * 1e-9
}
