package grade

import "almen-int/internal/coverage"

type Result struct {
	ArcHeight        float64
	Coverage         float64
	Plateau          float64
	Gain             float64
	Saturated        bool
	Recommend        Recommendation
	PeeningTime      float64
	SaturationRatio  float64
	CompleteCoverage bool
}

func Assemble(plateau, coverageValue, gain float64, sat coverage.Saturation) Result {
	arc := plateau * gain
	return Result{
		ArcHeight:        arc,
		Coverage:         coverageValue,
		Plateau:          plateau,
		Gain:             gain,
		Saturated:        sat.Saturated,
		Recommend:        Recommend(arc, sat.Saturated),
		PeeningTime:      sat.Time,
		SaturationRatio:  sat.Ratio,
		CompleteCoverage: sat.CompleteCoverage,
	}
}

func (r Result) GradeLetter() string {
	return r.Recommend.GradeLetter()
}

func (r Result) SaturationStateText() string {
	if r.Saturated {
		return "yes"
	}
	return "no"
}

func (r Result) HasGrade() bool {
	return r.Saturated && r.Recommend.Available
}
