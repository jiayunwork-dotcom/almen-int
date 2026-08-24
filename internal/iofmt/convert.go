package iofmt

import (
	"almen-int/internal/coverage"
	"almen-int/internal/model"
)

// BuildModelParams maps a validated case document onto the bending-model
// parameters. Callers must run Validate first; the conversion itself does not
// re-check the document.
func BuildModelParams(doc *CaseDoc) model.Params {
	return model.Params{
		Velocity:          *doc.Shot.Velocity,
		ReferenceVelocity: *doc.Reference.Velocity,
		ShotDiameter:      *doc.Shot.Diameter,
		ShotDensity:       *doc.Shot.Density,
		Thickness:         *doc.Strip.Thickness,
		Width:             *doc.Strip.Width,
		Length:            *doc.Strip.Length,
		Modulus:           *doc.Strip.Modulus,
		ResidualStress:    *doc.Residual.Stress,
		LayerDepth:        *doc.Residual.LayerDepth,
	}
}

// BuildCoverageParams maps a validated case document onto the coverage-model
// parameters.
func BuildCoverageParams(doc *CaseDoc) coverage.Params {
	return coverage.Params{
		Coverage:        *doc.Process.Coverage,
		RateConstant:    *doc.Process.RateConstant,
		GainCoefficient: *doc.Process.GainCoefficient,
	}
}

// ValueOrDefault dereferences an optional numeric field, returning the given
// fallback when the pointer is nil. It is used for the free-form metadata
// fields that do not feed the model.
func ValueOrDefault(v *float64, fallback float64) float64 {
	if v == nil {
		return fallback
	}
	return *v
}

// AllIssues runs the document-shape check and the two domain validators and
// returns every problem found, in a stable order: document, bending model,
// coverage. The domain validators only run when every required field is
// present, because the conversion would otherwise have nothing to dereference.
func AllIssues(doc *CaseDoc) []string {
	var issues []string
	issues = append(issues, Validate(doc)...)
	if len(issues) > 0 {
		return issues
	}
	issues = append(issues, model.Validate(BuildModelParams(doc))...)
	issues = append(issues, coverage.Validate(BuildCoverageParams(doc))...)
	return issues
}
