package model

type ModelReport struct {
	Params           Params
	VelocityRatio    float64
	ShotMassKG       float64
	KineticEnergy    float64
	ShotMomentum     float64
	EnergyDensity    float64
	ResidualStress   float64
	LayerForce       float64
	LeverArm         float64
	BendingMoment    float64
	MomentPerUnitW   float64
	MomentOfInertia  float64
	SectionModulus   float64
	ModulusPerArea   float64
	Curvature        float64
	ArcRadius        float64
	PlateauArcHeight float64
	ExactSagitta     float64
	RelativeError    float64
}

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
