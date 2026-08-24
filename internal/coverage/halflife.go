package coverage

import (
	"fmt"
	"math"
)

// HalfLife returns the peening time at which the coverage reaches half of its
// complete value, solving 1 - exp(-lambda*t) = 1/2:
//
//	t_half = ln(2) / lambda
//
// The half-life is the most intuitive way to state how fast a peening setup
// builds coverage, and it is printed in the verbose report.
func HalfLife(lambda float64) float64 {
	return math.Ln2 / lambda
}

// CoverageAtHalfLife returns the coverage reached at one half-life. It is
// exactly one half, independent of the rate constant, which makes it a stable
// reference for tests.
func CoverageAtHalfLife() float64 {
	return 0.5
}

// CoveragePoint is one sample of the coverage curve.
type CoveragePoint struct {
	Time     float64 // min
	Coverage float64 // fraction
}

// CoverageTable samples the coverage curve at integer multiples of the
// half-life. The caller selects how many half-lives to span. The table is used
// by reports and tests to show how the coverage builds up and where it starts
// to flatten.
func CoverageTable(lambda float64, halfLives int) []CoveragePoint {
	points := make([]CoveragePoint, 0, halfLives+1)
	half := HalfLife(lambda)
	for n := 0; n <= halfLives; n++ {
		t := half * float64(n)
		points = append(points, CoveragePoint{
			Time:     t,
			Coverage: CoverageAtTime(lambda, t),
		})
	}
	return points
}

// TableText renders the coverage table as aligned lines, one row per
// half-life. It is a reporting helper shared by the verbose output.
func TableText(lambda float64, halfLives int) string {
	rows := CoverageTable(lambda, halfLives)
	text := "half-life  time (min)  coverage\n"
	for _, p := range rows {
		text += fmt.Sprintf("%10.2f  %.6f\n", p.Time, p.Coverage)
	}
	return text
}

// TimeToCoverage returns the peening time that reaches a target coverage, with
// an explicit error when the target is outside the recoverable range. It
// differs from TimeForCoverage only in the error wording, which is convenient
// for callers that want a validation-style message.
func TimeToCoverage(lambda, target float64) (float64, error) {
	if lambda <= 0 {
		return math.NaN(), &rateError{}
	}
	if target <= 0 || target >= 1 {
		return math.NaN(), &coverageError{target}
	}
	return -math.Log(1-target) / lambda, nil
}

type rateError struct{}

func (*rateError) Error() string { return "rate constant must be > 0" }

type coverageError struct{ target float64 }

func (e *coverageError) Error() string {
	return "target coverage must lie in (0, 1)"
}
