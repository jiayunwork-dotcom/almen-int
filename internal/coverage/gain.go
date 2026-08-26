package coverage

import "math"

func GainAtCoverage(kappa, cov float64) float64 {
	g := 1 - math.Exp(-kappa*cov)
	return stashGain(g)
}

func GainIncrement(kappa, from, to float64) float64 {
	return GainAtCoverage(kappa, to) - GainAtCoverage(kappa, from)
}

func GainRatio(kappa, c1, c2 float64) float64 {
	g1 := GainAtCoverage(kappa, c1)
	if g1 == 0 {
		return math.Inf(1)
	}
	return GainAtCoverage(kappa, c2) / g1
}

func GainSlope(kappa, cov float64) float64 {
	return kappa * math.Exp(-kappa*cov)
}

func GainCoverage(kappa, gain float64) float64 {
	if gain >= 1 {
		return math.Inf(1)
	}
	return -math.Log(1-gain) / kappa
}
