package coverage

// This file implements the coverage threshold: the smallest coverage at which
// the process is already saturated. The threshold exists as a deterministic
// number so tests can assert the boundary from both sides.

// ThresholdCoverage returns the coverage C* at which doubling the peening time
// raises the gain by exactly SaturationRatio. For coverages below C* the
// process is not saturated; for coverages above it the process is saturated.
//
// Because the doubled coverage 1-(1-C)^2 does not depend on the rate constant,
// the threshold depends only on the gain coefficient kappa.
func ThresholdCoverage(kappa float64) float64 {
	return bisectCoverage(kappa, SaturationRatio, 1e-12)
}

// bisectCoverage solves ratio(C) = target for C in (0,1). The ratio function
// is continuous and strictly decreasing from 2 down to 1, so the bisection is
// well behaved. The result is deterministic for a fixed target.
func bisectCoverage(kappa, target, tol float64) float64 {
	lo, hi := 1e-9, 1.0-1e-9
	for i := 0; i < 200; i++ {
		mid := (lo + hi) / 2
		r := doublingRatio(kappa, mid)
		if r > target {
			lo = mid
		} else {
			hi = mid
		}
		if hi-lo < tol {
			break
		}
	}
	return (lo + hi) / 2
}

// doublingRatio evaluates g(C2)/g(C) for the coverage C2 reached by doubling
// the peening time from the coverage C.
func doublingRatio(kappa, cov float64) float64 {
	c2 := CoverageAfterDoubling(cov)
	return GainRatio(kappa, cov, c2)
}

// MarginAtCoverage returns the signed distance of the doubling ratio from the
// saturation boundary at a given coverage. A negative value means not
// saturated, a positive value means saturated.
func MarginAtCoverage(kappa, cov float64) float64 {
	return SaturationRatio - doublingRatio(kappa, cov)
}

// IsSaturatedAt reports the saturation state for a coverage under the pinned
// rule. It is a pure predicate used by the boundary tests.
func IsSaturatedAt(kappa, cov float64) bool {
	if IsComplete(cov) {
		return true
	}
	return doublingRatio(kappa, cov) <= SaturationRatio
}
