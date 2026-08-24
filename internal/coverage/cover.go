package coverage

import (
	"math"
)

// CoverageAtTime evaluates the exponential coverage law
//
//	C(t) = 1 - exp(-lambda * t)
//
// for a given rate constant and peening time. The value starts at zero, grows
// monotonically and approaches one as the time grows; it never exceeds one.
func CoverageAtTime(lambda, t float64) float64 {
	return 1 - math.Exp(-lambda*t)
}

// TimeForCoverage recovers the peening time that yields a given coverage:
//
//	t = -ln(1 - C) / lambda
//
// The inverse is only defined for coverages strictly below one, because a
// coverage of 1 is reached only in the limit of infinite time. For C >= 1 the
// function returns +Inf together with a nil error, which the saturation logic
// interprets as "complete coverage, nothing left to gain".
func TimeForCoverage(lambda, cov float64) (float64, error) {
	if cov >= 1 {
		return math.Inf(1), nil
	}
	if cov <= 0 {
		return math.Inf(-1), nil
	}
	return -math.Log(1-cov) / lambda, nil
}

// CoverageAfterDoubling returns the coverage reached when the peening time is
// doubled. Because C(t) = 1 - exp(-lambda*t), doubling the time maps a
// coverage C to
//
//	C2 = 1 - exp(-2*lambda*t) = 1 - (1-C)^2
//
// The mapping is independent of lambda, which is what makes the saturation
// threshold depend on the gain curve but not on the rate constant.
func CoverageAfterDoubling(cov float64) float64 {
	return 1 - (1-cov)*(1-cov)
}

// CoverageIncrement returns the additional coverage obtained by extending the
// peening time from t1 to t2. It is provided for reports that track how fast
// the coverage is still climbing.
func CoverageIncrement(lambda, t1, t2 float64) float64 {
	return CoverageAtTime(lambda, t2) - CoverageAtTime(lambda, t1)
}

// IsComplete reports whether the coverage has reached the complete mark
// (C >= 1), at which point the implied peening time is infinite and doubling
// it cannot change anything.
func IsComplete(cov float64) bool {
	return cov >= 1
}
