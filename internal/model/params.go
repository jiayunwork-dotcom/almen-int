package model

// velocityPower is the pinned exponent of the velocity power law. Shot kinetic
// energy scales with the square of the velocity and the induced residual
// stress follows the delivered energy density, so the exponent is pinned to 2.
// A bug that changes this exponent would break the "double the velocity,
// quadruple the arc height" invariant, which is exactly what the cross tests
// assert.
const velocityPower = 2.0

// Params carries every physical quantity the Almen bending model needs.
//
// Unit conventions are fixed and documented here because the model mixes them
// in one formula:
//
//	lengths and depths  : millimetres
//	stresses            : megapascals (N/mm^2)
//	elastic modulus     : gigapascals
//	velocities          : metres per second
//	densities           : kilograms per cubic metre
type Params struct {
	Velocity          float64 // shot velocity, m/s
	ReferenceVelocity float64 // velocity the quoted stress is valid at, m/s
	ShotDiameter      float64 // shot diameter, mm
	ShotDensity       float64 // shot material density, kg/m^3
	Thickness         float64 // strip thickness, mm
	Width             float64 // strip width, mm
	Length            float64 // strip span length, mm
	Modulus           float64 // strip elastic modulus, GPa
	ResidualStress    float64 // compressive stress at the reference velocity, MPa
	LayerDepth        float64 // depth of the residual compressive layer, mm
}

// VelocityRatio is the ratio v/v_ref of the shot velocity to the reference
// velocity at which the residual stress was quoted. It is the quantity that
// feeds the pinned power law and is exposed here so reports and tests agree on
// its definition.
func VelocityRatio(p Params) float64 {
	return p.Velocity / p.ReferenceVelocity
}

// VelocityExponent returns the pinned power-law exponent. The value is fixed
// to 2 so that a doubling of the velocity raises the kinetic energy, the
// residual stress and therefore the plateau arc height by a factor of four.
func VelocityExponent() float64 {
	return velocityPower
}

// ThicknessOf is a convenience accessor mirroring the strip thickness field.
// It exists so that sensitivity helpers and reports read the same source of
// truth when they reason about how thickness changes propagate.
func ThicknessOf(p Params) float64 {
	return p.Thickness
}

// ModulusOf is a convenience accessor for the elastic modulus in gigapascals.
func ModulusOf(p Params) float64 {
	return p.Modulus
}

// LayerDepthOf is a convenience accessor for the residual layer depth.
func LayerDepthOf(p Params) float64 {
	return p.LayerDepth
}
