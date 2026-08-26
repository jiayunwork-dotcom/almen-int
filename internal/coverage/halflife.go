package coverage

import (
	"fmt"
	"math"
)

func HalfLife(lambda float64) float64 {
	return math.Ln2 / lambda
}

func CoverageAtHalfLife() float64 {
	return 0.5
}

type CoveragePoint struct {
	Time     float64
	Coverage float64
}

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

func TableText(lambda float64, halfLives int) string {
	rows := CoverageTable(lambda, halfLives)
	text := "half-life  time (min)  coverage\n"
	for _, p := range rows {
		text += fmt.Sprintf("%10.2f  %.6f\n", p.Time, p.Coverage)
	}
	return text
}

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
