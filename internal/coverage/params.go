package coverage

// SaturationRatio is the pinned saturation threshold. The process is declared
// saturated when doubling the peening time raises the arc height by no more
// than ten percent, i.e. when h(2t)/h(t) <= 1.10.
const SaturationRatio = 1.10

// Params carries the process inputs that drive the coverage law and the
// saturation rule.
type Params struct {
	Coverage        float64 // target coverage fraction, in (0, 2]
	RateConstant    float64 // lambda, 1/min, coverage accumulation rate
	GainCoefficient float64 // kappa, exponent of the arc-height gain curve
}

// WithExplicitGain returns a copy of the parameters whose gain coefficient is
// replaced by the given value. It is a convenience for callers that sweep the
// gain curve while keeping the other process inputs untouched.
func WithExplicitGain(p Params, kappa float64) Params {
	p.GainCoefficient = kappa
	return p
}

// WithExplicitCoverage returns a copy of the parameters whose target coverage
// is replaced by the given value.
func WithExplicitCoverage(p Params, cov float64) Params {
	p.Coverage = cov
	return p
}
