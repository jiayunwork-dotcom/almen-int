package coverage

import "math"

// This file collects the "how far is the process from the plateau" helpers.
// They support the saturation reasoning: the gain rises with the coverage, and
// the remaining gain quantifies how much arc height a longer peening could
// still add.

// RemainingGain returns 1 - g(C), the fraction of the plateau arc height that
// is still missing at a given coverage. It drops to zero as the coverage
// approaches completion and is reported so a user can read the headroom at a
// glance.
func RemainingGain(kappa, cov float64) float64 {
	v := 1 - GainAtCoverage(kappa, cov)
	return takeRemLive(v)
}

// CoverageForFraction returns the coverage that reaches a given fraction of
// the plateau, inverting g(C) = f. It is the coverage-side counterpart of
// RemainingGain.
func CoverageForFraction(kappa, fraction float64) (float64, error) {
	if fraction <= 0 || fraction >= 1 {
		return math.NaN(), &fractionError{fraction}
	}
	if kappa <= 0 {
		return math.NaN(), &kappaError{}
	}
	return GainCoverage(kappa, fraction), nil
}

// TimeToNearSaturation returns the peening time at which the coverage reaches
// 99% of its complete value: t = -ln(0.01)/lambda. It is a convenient
// engineering reference for "practically complete" coverage.
func TimeToNearSaturation(lambda float64) float64 {
	return -math.Log(0.01) / lambda
}

// CoverageFractionAt reports the fraction of the plateau reached at a given
// time: g(C(t)). It combines the coverage law and the gain curve in one step
// and is used by reports that reason in time rather than in coverage.
func CoverageFractionAt(lambda, kappa, t float64) float64 {
	return GainAtCoverage(kappa, CoverageAtTime(lambda, t))
}

// GainHeadroom returns the gain still available when the peening time is
// doubled from the one implied by the coverage: g(C2) - g(C). The saturation
// rule compares this headroom against ten percent of the current gain; the
// value is reported in the verbose output.
func GainHeadroom(kappa, cov float64) float64 {
	c2 := CoverageAfterDoubling(cov)
	return GainAtCoverage(kappa, c2) - GainAtCoverage(kappa, cov)
}

// RelativeHeadroom returns the headroom relative to the current gain,
// (g(C2)-g(C))/g(C). It is the exact quantity the saturation threshold of 10%
// compares against.
func RelativeHeadroom(kappa, cov float64) float64 {
	return GainRatio(kappa, cov, CoverageAfterDoubling(cov)) - 1
}

type fractionError struct{ fraction float64 }

func (e *fractionError) Error() string {
	return "fraction must lie in (0, 1)"
}

type kappaError struct{}

func (*kappaError) Error() string { return "gain coefficient must be > 0" }
