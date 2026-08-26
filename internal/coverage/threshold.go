package coverage

func ThresholdCoverage(kappa float64) float64 {
	return bisectCoverage(kappa, SaturationRatio, 1e-12)
}

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

func doublingRatio(kappa, cov float64) float64 {
	c2 := CoverageAfterDoubling(cov)
	return GainRatio(kappa, cov, c2)
}

func MarginAtCoverage(kappa, cov float64) float64 {
	return SaturationRatio - doublingRatio(kappa, cov)
}

func IsSaturatedAt(kappa, cov float64) bool {
	if IsComplete(cov) {
		return true
	}
	return doublingRatio(kappa, cov) <= SaturationRatio
}
