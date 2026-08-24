package coverage

import (
	"math"
)

// Saturation is the outcome of applying the pinned saturation rule to a set of
// process inputs.
type Saturation struct {
	Params           Params
	Time             float64 // implied peening time for the target coverage, min
	DoubledTime      float64 // twice the implied time, min
	DoubledCoverage  float64 // coverage reached after doubling the time
	Gain             float64 // g(C) at the target coverage
	DoubledGain      float64 // g(C2) at the doubled coverage
	Ratio            float64 // g(C2)/g(C), the arc-height growth from doubling
	Threshold        float64 // the ratio that would sit exactly at the boundary
	Saturated        bool    // ratio <= SaturationRatio
	CompleteCoverage bool    // C >= 1, doubling time leaves nothing to gain
}

// Determine applies the saturation rule to the process inputs. It validates
// the parameters first and returns the list of violations as an error when
// any rule fails.
//
// For a coverage below one the implied time is recovered, doubled, and the
// gain ratio is compared against the pinned threshold. For a coverage at or
// above one the process is complete and therefore saturated by definition:
// the implied time is already infinite and doubling it changes nothing.
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

// Margin returns how far the doubling ratio sits from the saturation boundary.
// A negative margin means the process still has more than ten percent growth
// left and is not saturated; a positive margin means it is past the boundary.
func (s Saturation) Margin() float64 {
	return s.Threshold - s.Ratio
}

// GrowthPercent returns the relative growth of the arc height when the peening
// time is doubled, as a percentage. It is the quantity the CLI prints so the
// saturation decision can be read directly.
func (s Saturation) GrowthPercent() float64 {
	return (s.Ratio - 1) * 100
}

// issuesError turns a list of validation messages into an error. The messages
// are joined on newlines so the CLI can print them verbatim on stderr.
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

// validationError is a trivial error wrapper that preserves the exact message
// text produced by the validation rules.
type validationError struct{ msg string }

func (e *validationError) Error() string { return e.msg }
