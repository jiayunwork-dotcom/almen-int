package model

// ModelReport is a snapshot of every intermediate quantity produced by the
// bending model for one parameter set. Keeping all of them in one struct makes
// the verbose CLI output and the JSON output share a single source of truth.
type ModelReport struct {
	Params           Params
	VelocityRatio    float64 // v / v_ref
	ShotMassKG       float64 // mass of one shot, kg
	KineticEnergy    float64 // energy of one shot, J
	ShotMomentum     float64 // momentum of one shot, kg*m/s
	EnergyDensity    float64 // energy per unit shot volume, J/mm^3
	ResidualStress   float64 // MPa, at the case velocity
	LayerForce       float64 // N
	LeverArm         float64 // mm
	BendingMoment    float64 // N*mm
	MomentPerUnitW   float64 // N
	MomentOfInertia  float64 // mm^4
	SectionModulus   float64 // mm^3
	ModulusPerArea   float64 // N/mm^2
	Curvature        float64 // 1/mm
	ArcRadius        float64 // mm
	PlateauArcHeight float64 // mm, fully peened
	ExactSagitta     float64 // mm
	RelativeError    float64 // |exact - h| / exact
}

// BuildReport evaluates the full bending chain for a parameter set and returns
// the snapshot. Callers are expected to validate the parameters before calling
// this; the report does not re-check them.
func BuildReport(p Params) ModelReport {
	return ModelReport{
		Params:           p,
		VelocityRatio:    VelocityRatio(p),
		ShotMassKG:       ShotMass(p),
		KineticEnergy:    KineticEnergy(p),
		ShotMomentum:     ShotMomentum(p),
		EnergyDensity:    EnergyDensity(p),
		ResidualStress:   ResidualStress(p),
		LayerForce:       LayerForce(p),
		LeverArm:         LeverArm(p),
		BendingMoment:    BendingMoment(p),
		MomentPerUnitW:   MomentPerUnitWidth(p),
		MomentOfInertia:  MomentOfInertia(p),
		SectionModulus:   SectionModulus(p),
		ModulusPerArea:   ModulusPerArea(p),
		Curvature:        Curvature(p),
		ArcRadius:        ArcRadius(p),
		PlateauArcHeight: PlateauArcHeight(p),
		ExactSagitta:     ExactSagitta(p),
		RelativeError:    RelativeError(p),
	}
}
