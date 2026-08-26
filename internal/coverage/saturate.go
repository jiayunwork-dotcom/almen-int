package coverage

import (
	"math"
)

type Saturation struct {
	Params           Params
	Time             float64
	DoubledTime      float64
	DoubledCoverage  float64
	Gain             float64
	DoubledGain      float64
	Ratio            float64
	Threshold        float64
	Saturated        bool
	CompleteCoverage bool
}

func Determine(p Params) (Saturation, error) {
	if issues := Validate(p); len(issues) > 0 {
		return Saturation{}, issuesError(issues)
	}

	sat := Saturation{
		Params:    p,
		Threshold: SaturationRatio,
	}

	gain := GainAtCoverage(p.GainCoefficient, p.Coverage)
	sat.Gain = gain

	if IsComplete(p.Coverage) {
		sat.Time = math.Inf(1)
		sat.DoubledTime = math.Inf(1)
		sat.DoubledCoverage = 1
		sat.DoubledGain = 1
		sat.Ratio = 1
		sat.CompleteCoverage = true
		sat.Saturated = true
		return sat, nil
	}

	t, err := TimeForCoverage(p.RateConstant, p.Coverage)
	if err != nil || math.IsInf(t, 0) {
		return Saturation{}, issuesError([]string{"cannot recover a finite peening time for the target coverage"})
	}
	sat.Time = t
	sat.DoubledTime = 2 * t

	c2 := CoverageAfterDoubling(p.Coverage)
	sat.DoubledCoverage = c2
	sat.DoubledGain = GainAtCoverage(p.GainCoefficient, c2)
	sat.Ratio = sat.DoubledGain / gain
	sat.Saturated = sat.Ratio <= SaturationRatio
	return sat, nil
}

func (s Saturation) Margin() float64 {
	return s.Threshold - s.Ratio
}

func (s Saturation) GrowthPercent() float64 {
	return (s.Ratio - 1) * 100
}

func issuesError(issues []string) error {
	var text string
	for i, msg := range issues {
		if i > 0 {
			text += "\n"
		}
		text += msg
	}
	return &validationError{text}
}

type validationError struct{ msg string }

func (e *validationError) Error() string { return e.msg }
