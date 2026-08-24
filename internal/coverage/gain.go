package coverage

import "math"

// GainAtCoverage returns the arc-height gain g(C) = 1 - exp(-kappa*C) that
// scales the fully peened plateau arc height down to the value reached at a
// partial coverage. The gain grows monotonically from zero, is concave, and
// approaches one as the coverage grows:
//
//	g(0) = 0        no peening, no arc height
//	g(C) -> 1       complete coverage approaches the plateau
//
// The gain is the only place where the coverage influences the arc height.
// The bending model itself never sees the coverage, which is a deliberate
// separation: a bug that routed the coverage into the moment formula would
// break the invariant tested by the cross-package tests.
func GainAtCoverage(kappa, cov float64) float64 {
	return 1 - math.Exp(-kappa*cov)
}

// GainIncrement returns g(to) - g(from). Because the gain is concave, the
// increment between 0.5 and 1.0 is always smaller than the increment between
// 0 and 0.5, which is the coverage cross rule the tests assert.
func GainIncrement(kappa, from, to float64) float64 {
	return GainAtCoverage(kappa, to) - GainAtCoverage(kappa, from)
}

// GainRatio returns g(C2)/g(C1), the relative growth of the gain when the
// coverage moves from C1 to C2. The saturation rule compares exactly this
// ratio for the coverage reached by doubling the peening time.
func GainRatio(kappa, c1, c2 float64) float64 {
	g1 := GainAtCoverage(kappa, c1)
	if g1 == 0 {
		return math.Inf(1)
	}
	return GainAtCoverage(kappa, c2) / g1
}

// GainSlope returns the derivative kappa*exp(-kappa*C) of the gain curve. The
// slope is monotonically decreasing, which quantifies the diminishing returns
// of additional coverage.
func GainSlope(kappa, cov float64) float64 {
	return kappa * math.Exp(-kappa*cov)
}

// GainCoverage reports, for a given gain value, the coverage that reaches it,
// solving C = -ln(1-g)/kappa. It is the inverse of GainAtCoverage and is used
// by the reports to talk about the process in coverage terms.
func GainCoverage(kappa, gain float64) float64 {
	if gain >= 1 {
		return math.Inf(1)
	}
	return -math.Log(1-gain) / kappa
}
