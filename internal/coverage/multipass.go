package coverage

import (
	"errors"
	"fmt"
	"math"
)

func CombineIndependent(parts []float64) (float64, error) {
	if len(parts) == 0 {
		return 0, errors.New("coverage: at least one pass is required")
	}
	uncovered := 1.0
	for i, c := range parts {
		if c < 0 || math.IsNaN(c) || math.IsInf(c, 0) {
			return 0, fmt.Errorf("coverage: pass %d coverage must be finite and non-negative", i)
		}
		frac := c
		if frac > 1 {
			frac = 1
		}
		uncovered *= 1 - frac
	}
	out := 1 - uncovered
	if out < 0 {
		return 0, nil
	}
	if out > 1 {
		return 1, nil
	}
	return out, nil
}

func CombineTimes(lambda float64, times []float64) (float64, error) {
	if lambda <= 0 || math.IsNaN(lambda) || math.IsInf(lambda, 0) {
		return 0, errors.New("coverage: rate constant must be positive")
	}
	sum := 0.0
	for i, t := range times {
		if t < 0 || math.IsNaN(t) || math.IsInf(t, 0) {
			return 0, fmt.Errorf("coverage: pass %d time must be finite and non-negative", i)
		}
		sum += t
	}
	return CoverageAtTime(lambda, sum), nil
}

func CombineIndependentAgreesWithTimes(lambda float64, times []float64) (float64, error) {
	parts := make([]float64, len(times))
	for i, t := range times {
		parts[i] = CoverageAtTime(lambda, t)
	}
	fromParts, err := CombineIndependent(parts)
	if err != nil {
		return 0, err
	}
	fromTimes, err := CombineTimes(lambda, times)
	if err != nil {
		return 0, err
	}
	if math.Abs(fromParts-fromTimes) > 1e-12 {
		return 0, fmt.Errorf("coverage: Avrami product %g disagrees with time-sum %g", fromParts, fromTimes)
	}
	return fromTimes, nil
}

func RemainingUncovered(cov float64) float64 {
	if cov >= 1 {
		return 0
	}
	if cov <= 0 {
		return 1
	}
	return 1 - cov
}

func OrderIndependent(a, b []float64) (bool, error) {
	x, err := CombineIndependent(a)
	if err != nil {
		return false, err
	}
	y, err := CombineIndependent(b)
	if err != nil {
		return false, err
	}
	return math.Abs(x-y) <= 1e-12, nil
}
