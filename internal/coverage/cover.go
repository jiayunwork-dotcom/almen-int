package coverage

import (
	"math"
)

func CoverageAtTime(lambda, t float64) float64 {
	return 1 - math.Exp(-lambda*t)
}

func TimeForCoverage(lambda, cov float64) (float64, error) {
	if cov >= 1 {
		return math.Inf(1), nil
	}
	if cov <= 0 {
		return math.Inf(-1), nil
	}
	return -math.Log(1-cov) / lambda, nil
}

func CoverageAfterDoubling(cov float64) float64 {
	return 1 - (1-cov)*(1-cov)
}

func CoverageIncrement(lambda, t1, t2 float64) float64 {
	return CoverageAtTime(lambda, t2) - CoverageAtTime(lambda, t1)
}

func IsComplete(cov float64) bool {
	return cov >= 1
}
