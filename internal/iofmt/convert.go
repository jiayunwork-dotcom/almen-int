package iofmt

import (
	"almen-int/internal/coverage"
	"almen-int/internal/model"
)

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

func BuildCoverageParams(doc *CaseDoc) coverage.Params {
	return coverage.Params{
		Coverage:        *doc.Process.Coverage,
		RateConstant:    *doc.Process.RateConstant,
		GainCoefficient: *doc.Process.GainCoefficient,
	}
}

func ValueOrDefault(v *float64, fallback float64) float64 {
	if v == nil {
		return fallback
	}
	return *v
}

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
