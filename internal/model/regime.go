package model

type DeflectionRegime struct {
	ThicknessOverSpan float64
	ArcHeightOverSpan float64
	RadiusOverSpan    float64
	SurfaceStrain     float64
	PeakBendingStress float64
	ThinStripOK       bool
	SmallDeflectionOK bool
	ElasticOK         bool
}

const (
	ThinStripLimit       = 0.05
	SmallDeflectionLimit = 0.01
)

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

func (r DeflectionRegime) RegimeSummary() string {
	if !r.ThinStripOK {
		return "strip is too thick relative to its span for the thin-strip model"
	}
	if !r.SmallDeflectionOK {
		return "deflection is too large for the small-deflection formula"
	}
	return "within the thin-strip and small-deflection regime"
}
