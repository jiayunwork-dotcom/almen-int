package coverage

const SaturationRatio = 1.10

type Params struct {
	Coverage        float64
	RateConstant    float64
	GainCoefficient float64
}

func WithExplicitGain(p Params, kappa float64) Params {
	p.GainCoefficient = kappa
	return p
}

func WithExplicitCoverage(p Params, cov float64) Params {
	p.Coverage = cov
	return p
}
