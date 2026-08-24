package model

// DeflectionRegime quantifies how well the case fits the thin-strip and
// small-deflection assumptions that justify the pinned arc-height formula.
// Reporting these dimensionless numbers is useful because the model is only
// calibrated for strips that actually behave like thin beams.
type DeflectionRegime struct {
	ThicknessOverSpan float64 // t / L
	ArcHeightOverSpan float64 // h / L
	RadiusOverSpan    float64 // R / L
	SurfaceStrain     float64 // t / (2R), the peak elastic strain at the surfaces
	PeakBendingStress float64 // E * t * kappa / 2, N/mm^2
	ThinStripOK       bool    // t / L is small enough
	SmallDeflectionOK bool    // h / L is small enough
	ElasticOK         bool    // the peak bending stress stays below the residual stress
}

// ThinStripLimit and SmallDeflectionLimit are the pinned thresholds used by
// the regime checks. They are generous on purpose: real Almen strips sit far
// below both limits, so the checks only fire for clearly unreasonable input.
const (
	ThinStripLimit       = 0.05 // t / L below 5%
	SmallDeflectionLimit = 0.01 // h / L below 1%
)

// EvaluateRegime computes the dimensionless checks for a parameter set. The
// checks do not reject anything; they only annotate the report with how far
// the case is from the model's assumptions.
func EvaluateRegime(p Params) DeflectionRegime {
	kappa := Curvature(p)
	h := PlateauArcHeight(p)
	R := 1 / kappa
	strain := p.Thickness / (2 * R)
	peakStress := ModulusPerArea(p) * p.Thickness * kappa / 2

	tOverL := p.Thickness / p.Length
	hOverL := h / p.Length
	return DeflectionRegime{
		ThicknessOverSpan: tOverL,
		ArcHeightOverSpan: hOverL,
		RadiusOverSpan:    R / p.Length,
		SurfaceStrain:     strain,
		PeakBendingStress: peakStress,
		ThinStripOK:       tOverL <= ThinStripLimit,
		SmallDeflectionOK: hOverL <= SmallDeflectionLimit,
		ElasticOK:         peakStress < p.ResidualStress,
	}
}

// RegimeSummary returns a one-line human description of the regime verdict.
func (r DeflectionRegime) RegimeSummary() string {
	if !r.ThinStripOK {
		return "strip is too thick relative to its span for the thin-strip model"
	}
	if !r.SmallDeflectionOK {
		return "deflection is too large for the small-deflection formula"
	}
	return "within the thin-strip and small-deflection regime"
}
