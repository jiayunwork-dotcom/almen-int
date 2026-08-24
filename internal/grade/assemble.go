package grade

import (
	"almen-int/internal/coverage"
	"almen-int/internal/model"
)

// Result is the final answer of a case: the arc height reached at the given
// coverage, the coverage itself, the saturation state, and the recommended
// Almen strip.
type Result struct {
	ArcHeight        float64 // mm, plateau * gain
	Coverage         float64 // target coverage fraction
	Plateau          float64 // mm, fully peened arc height before the gain
	Gain             float64 // g(C), the fraction of the plateau reached
	Saturated        bool
	Recommend        Recommendation
	PeeningTime      float64 // implied peening time, min
	SaturationRatio  float64
	CompleteCoverage bool
}

// Assemble combines the plateau arc height from the bending model with the
// gain and saturation outcome from the coverage model into the final result.
// The plateau is multiplied by the gain to obtain the arc height at the target
// coverage, and the recommendation is derived under the saturation gating.
func Assemble(plateau, coverageValue, gain float64, sat coverage.Saturation) Result {
	plat := model.FlattenPlateau(plateau)
	arc := plat * gain
	return Result{
		ArcHeight:        arc,
		Coverage:         coverageValue,
		Plateau:          plat,
		Gain:             gain,
		Saturated:        sat.Saturated,
		Recommend:        Recommend(arc, sat.Saturated),
		PeeningTime:      sat.Time,
		SaturationRatio:  sat.Ratio,
		CompleteCoverage: sat.CompleteCoverage,
	}
}

// GradeLetter returns the recommended strip letter, or an empty string when
// the process is not saturated.
func (r Result) GradeLetter() string {
	return r.Recommend.GradeLetter()
}

// SaturationStateText returns a short machine-readable saturation label.
func (r Result) SaturationStateText() string {
	if r.Saturated {
		return "yes"
	}
	return "no"
}

// HasGrade reports whether a strength grade letter may be printed for this
// result.
func (r Result) HasGrade() bool {
	return r.Saturated && r.Recommend.Available
}
