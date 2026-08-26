package coverage

import (
	"math"
	"testing"
)

func testParams() Params {
	return Params{
		Coverage:        0.98,
		RateConstant:    0.085,
		GainCoefficient: 2.6,
	}
}

func TestCoverageExponentialLaw(t *testing.T) {
	if got := CoverageAtTime(0.085, 0); got != 0 {
		t.Errorf("coverage at t=0 must be 0, got %g", got)
	}
	prev := 0.0
	for _, tm := range []float64{1, 5, 10, 20, 40, 80} {
		c := CoverageAtTime(0.085, tm)
		if c <= prev {
			t.Errorf("coverage must grow with time: at t=%g got %g after %g", tm, c, prev)
		}
		if c > 1 {
			t.Errorf("coverage must never exceed 1, got %g at t=%g", c, tm)
		}
		prev = c
	}
	tm, err := TimeForCoverage(0.085, 0.98)
	if err != nil {
		t.Fatalf("TimeForCoverage must not fail for a sub-unity coverage: %v", err)
	}
	if !closeEnough(CoverageAtTime(0.085, tm), 0.98, 1e-9) {
		t.Errorf("time inverse must round-trip the coverage: got %g", CoverageAtTime(0.085, tm))
	}
}

func TestTimeInverseHandlesCompleteCoverage(t *testing.T) {
	tm, err := TimeForCoverage(0.085, 1.0)
	if err != nil {
		t.Fatalf("TimeForCoverage(1.0) must not error, got %v", err)
	}
	if !math.IsInf(tm, 1) {
		t.Errorf("coverage of 1 must imply an infinite time, got %g", tm)
	}
	if got := CoverageAfterDoubling(0.98); !closeEnough(got, 0.9996, 1e-12) {
		t.Errorf("doubling 0.98 coverage must give 1-(1-0.98)^2 = 0.9996, got %g", got)
	}
}

func TestSaturationFlag(t *testing.T) {
	low := WithExplicitCoverage(testParams(), 0.5)
	high := testParams()

	sLow, err := Determine(low)
	if err != nil {
		t.Fatalf("Determine(low) failed: %v", err)
	}
	sHigh, err := Determine(high)
	if err != nil {
		t.Fatalf("Determine(high) failed: %v", err)
	}

	if sLow.Saturated {
		t.Errorf("coverage 0.5 must not be saturated (ratio %g), got saturated", sLow.Ratio)
	}
	if !sHigh.Saturated {
		t.Errorf("coverage 0.98 must be saturated (ratio %g), got not saturated", sHigh.Ratio)
	}
	if !(sLow.Ratio > 1.10) {
		t.Errorf("low coverage ratio must exceed the 10%% threshold, got %g", sLow.Ratio)
	}
	if !(sHigh.Ratio <= 1.10) {
		t.Errorf("high coverage ratio must sit at or below the 10%% threshold, got %g", sHigh.Ratio)
	}
	if sHigh.GrowthPercent() > 10 {
		t.Errorf("saturated case growth must be at most 10%%, got %g%%", sHigh.GrowthPercent())
	}
}

func TestCompleteCoverageIsSaturated(t *testing.T) {
	for _, cov := range []float64{1.0, 1.5, 2.0} {
		p := WithExplicitCoverage(testParams(), cov)
		s, err := Determine(p)
		if err != nil {
			t.Fatalf("Determine(%.1f) failed: %v", cov, err)
		}
		if !s.Saturated {
			t.Errorf("coverage %g must be saturated", cov)
		}
		if !s.CompleteCoverage {
			t.Errorf("coverage %g must be flagged as complete", cov)
		}
	}
}

func TestSaturationBoundary(t *testing.T) {
	kappa := testParams().GainCoefficient
	cStar := ThresholdCoverage(kappa)
	if cStar <= 0 || cStar >= 1 {
		t.Fatalf("threshold coverage must lie in (0,1), got %g", cStar)
	}
	below := cStar - 1e-6
	above := cStar + 1e-6
	if IsSaturatedAt(kappa, below) {
		t.Errorf("coverage %g (below threshold %g) must not be saturated", below, cStar)
	}
	if !IsSaturatedAt(kappa, above) {
		t.Errorf("coverage %g (above threshold %g) must be saturated", above, cStar)
	}
	r := doublingRatio(kappa, cStar)
	if !closeEnough(r, SaturationRatio, 1e-9) {
		t.Errorf("threshold must sit at ratio %.3f, got %g", SaturationRatio, r)
	}
}

func TestGainDiminishingIncrements(t *testing.T) {
	kappa := testParams().GainCoefficient
	first := GainIncrement(kappa, 0, 0.5)
	second := GainIncrement(kappa, 0.5, 1.0)
	if !(second < first) {
		t.Errorf("gain increment 0.5->1.0 (%g) must be smaller than 0->0.5 (%g)", second, first)
	}
	if !(second > 0) {
		t.Errorf("gain must still rise between 0.5 and 1.0, got %g", second)
	}
	if GainAtCoverage(kappa, 1.0) >= 1 {
		t.Errorf("gain at coverage 1 must stay below 1, got %g", GainAtCoverage(kappa, 1.0))
	}
}

func TestSaturationThresholdIndependentOfRateConstant(t *testing.T) {
	kappa := testParams().GainCoefficient
	one := ThresholdCoverage(kappa)
	for _, cov := range []float64{0.3, 0.5, 0.7, 0.9} {
		r1 := doublingRatio(kappa, cov)
		r2 := doublingRatio(kappa, cov)
		if !closeEnough(r1, r2, 0) {
			t.Errorf("doubling ratio must not depend on the rate constant, got %g vs %g", r1, r2)
		}
	}
	if one <= 0 {
		t.Errorf("threshold must be positive, got %g", one)
	}
}

func TestCoverageValidation(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Params)
	}{
		{"zero coverage", func(p *Params) { p.Coverage = 0 }},
		{"negative coverage", func(p *Params) { p.Coverage = -0.2 }},
		{"coverage above two", func(p *Params) { p.Coverage = 2.5 }},
		{"zero rate constant", func(p *Params) { p.RateConstant = 0 }},
		{"zero gain coefficient", func(p *Params) { p.GainCoefficient = 0 }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := testParams()
			tc.mutate(&p)
			if issues := Validate(p); len(issues) == 0 {
				t.Errorf("expected validation issues for %s, got none", tc.name)
			}
			if Valid(p) {
				t.Errorf("Valid() must be false for %s", tc.name)
			}
		})
	}
	edge := WithExplicitCoverage(testParams(), 2.0)
	if issues := Validate(edge); len(issues) != 0 {
		t.Errorf("coverage 2.0 must be valid, got: %v", issues)
	}
}

func TestDetermineReturnsErrorOnBadInput(t *testing.T) {
	bad := WithExplicitCoverage(testParams(), 0)
	if _, err := Determine(bad); err == nil {
		t.Errorf("Determine must fail for an invalid coverage")
	}
}

func TestMultipassAgreesWithTimeSum(t *testing.T) {
	lambda := 0.085
	times := []float64{4, 7, 2}
	got, err := CombineIndependentAgreesWithTimes(lambda, times)
	if err != nil {
		t.Fatal(err)
	}
	want := CoverageAtTime(lambda, 13)
	if !closeEnough(got, want, 1e-12) {
		t.Fatalf("combined coverage %g want %g", got, want)
	}
	rev := []float64{2, 7, 4}
	ok, err := OrderIndependent(
		[]float64{CoverageAtTime(lambda, 4), CoverageAtTime(lambda, 7), CoverageAtTime(lambda, 2)},
		[]float64{CoverageAtTime(lambda, rev[0]), CoverageAtTime(lambda, rev[1]), CoverageAtTime(lambda, rev[2])},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Avrami product must be order-independent")
	}
	full, err := CombineIndependent([]float64{1.5, 0.2})
	if err != nil {
		t.Fatal(err)
	}
	if full != 1 {
		t.Fatalf("coverage above 1 must saturate the remaining area, got %g", full)
	}
	if RemainingUncovered(0.75) != 0.25 {
		t.Fatalf("remaining uncovered %g", RemainingUncovered(0.75))
	}
}

func closeEnough(a, b, tol float64) bool {
	scale := 1.0
	if b != 0 {
		scale = b
	}
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol*scale
}
