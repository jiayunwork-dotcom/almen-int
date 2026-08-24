package model

// This file holds sensitivity helpers that quantify how a change in a single
// process variable propagates to the plateau arc height. They exist so that
// the cross rules from the project contract can be stated as code and asserted
// by tests: double the velocity, double the thickness, and so on.

// VelocityScaling returns the ratio h(hi)/h(lo) of the plateau arc heights
// evaluated at two velocities, everything else held at the "hi" parameter set.
// With the pinned exponent of 2 the ratio is exactly (v_hi/v_lo)^2, so a
// velocity doubling yields a factor of four.
func VelocityScaling(lo, hi float64, p Params) float64 {
	slow := p
	slow.Velocity = lo
	fast := p
	fast.Velocity = hi
	return PlateauArcHeight(fast) / PlateauArcHeight(slow)
}

// ThicknessScaling returns the ratio h(t_hi)/h(t_lo) of the plateau arc
// heights evaluated at two strip thicknesses, everything else held at the
// given parameter set. Because the inertia is cubic in the thickness while the
// moment only grows linearly with the lever arm, the ratio is below one for
// any realistic thickness pair: a thicker strip bends less.
func ThicknessScaling(lo, hi float64, p Params) float64 {
	thin := p
	thin.Thickness = lo
	thick := p
	thick.Thickness = hi
	return PlateauArcHeight(thick) / PlateauArcHeight(thin)
}

// InertiaScaling returns the ratio of the second moments of area for two
// thicknesses: I(t_hi)/I(t_lo) = (t_hi/t_lo)^3. A doubling of the thickness
// raises the inertia by exactly eight.
func InertiaScaling(lo, hi float64) float64 {
	return MomentOfInertiaFor(1, hi) / MomentOfInertiaFor(1, lo)
}

// StressScaling returns sigma(hi)/sigma(lo) for two velocities under the
// pinned power law: (v_hi/v_lo)^2. It is provided separately from
// VelocityScaling so a report can show that the stress and the arc height
// follow the same power while the kinetic energy does so through the shot
// mass as well.
func StressScaling(lo, hi float64, p Params) float64 {
	slow := p
	slow.Velocity = lo
	fast := p
	fast.Velocity = hi
	return ResidualStress(fast) / ResidualStress(slow)
}

// EnergyScaling returns the ratio of the shot kinetic energies at two
// velocities: KE(hi)/KE(lo) = (v_hi/v_lo)^2. The kinetic energy carries the
// same exponent as the stress, which is the physical justification for the
// pinned power law.
func EnergyScaling(lo, hi float64, p Params) float64 {
	slow := p
	slow.Velocity = lo
	fast := p
	fast.Velocity = hi
	return KineticEnergy(fast) / KineticEnergy(slow)
}
