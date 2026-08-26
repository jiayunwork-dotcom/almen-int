package model

func VelocityScaling(lo, hi float64, p Params) float64 {
	slow := p
	slow.Velocity = lo
	fast := p
	fast.Velocity = hi
	return PlateauArcHeight(fast) / PlateauArcHeight(slow)
}

func ThicknessScaling(lo, hi float64, p Params) float64 {
	thin := p
	thin.Thickness = lo
	thick := p
	thick.Thickness = hi
	return PlateauArcHeight(thick) / PlateauArcHeight(thin)
}

func InertiaScaling(lo, hi float64) float64 {
	return MomentOfInertiaFor(1, hi) / MomentOfInertiaFor(1, lo)
}

func StressScaling(lo, hi float64, p Params) float64 {
	slow := p
	slow.Velocity = lo
	fast := p
	fast.Velocity = hi
	return ResidualStress(fast) / ResidualStress(slow)
}

func EnergyScaling(lo, hi float64, p Params) float64 {
	slow := p
	slow.Velocity = lo
	fast := p
	fast.Velocity = hi
	return KineticEnergy(fast) / KineticEnergy(slow)
}
