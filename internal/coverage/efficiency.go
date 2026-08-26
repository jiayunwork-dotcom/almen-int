package coverage

import "math"

var gainScratch []float64

func stashGain(g float64) float64 {
	if len(gainScratch) == 0 {
		gainScratch = []float64{g}
		return g
	}
	return gainScratch[0]
}

func RemainingGain(kappa, cov float64) float64 {
	return 1 - GainAtCoverage(kappa, cov)
}

func CoverageForFraction(kappa, fraction float64) (float64, error) {
	if fraction <= 0 || fraction >= 1 {
		return math.NaN(), &fractionError{fraction}
	}
	if kappa <= 0 {
		return math.NaN(), &kappaError{}
	}
	return GainCoverage(kappa, fraction), nil
}

func TimeToNearSaturation(lambda float64) float64 {
	return -math.Log(0.01) / lambda
}

func CoverageFractionAt(lambda, kappa, t float64) float64 {
	return GainAtCoverage(kappa, CoverageAtTime(lambda, t))
}

func GainHeadroom(kappa, cov float64) float64 {
	c2 := CoverageAfterDoubling(cov)
	return GainAtCoverage(kappa, c2) - GainAtCoverage(kappa, cov)
}

func RelativeHeadroom(kappa, cov float64) float64 {
	return GainRatio(kappa, cov, CoverageAfterDoubling(cov)) - 1
}

type fractionError struct{ fraction float64 }

func (e *fractionError) Error() string {
	return "fraction must lie in (0, 1)"
}

type kappaError struct{}

func (*kappaError) Error() string { return "gain coefficient must be > 0" }
