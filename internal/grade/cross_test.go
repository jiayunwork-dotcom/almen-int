package grade

import (
	"testing"

	"almen-int/internal/coverage"
	"almen-int/internal/model"
)

// This file holds cross-package tests that bind the bending model and the
// coverage model together. They assert the cross rules of the project
// contract: velocity only moves the arc height, thickness only moves it the
// other way, and the saturation state never depends on either.

func testModelParams() model.Params {
	return model.Params{
		Velocity:          48.0,
		ReferenceVelocity: 50.0,
		ShotDiameter:      0.6,
		ShotDensity:       7800.0,
		Thickness:         1.29,
		Width:             18.5,
		Length:            76.0,
		Modulus:           205.0,
		ResidualStress:    850.0,
		LayerDepth:        0.05,
	}
}

func fullResult(v float64, t float64, cov float64, kappa float64) Result {
	mp := testModelParams()
	mp.Velocity = v
	mp.Thickness = t
	cp := coverage.Params{
		Coverage:        cov,
		RateConstant:    0.085,
		GainCoefficient: kappa,
	}
	plateau := model.PlateauArcHeight(mp)
	gain := coverage.GainAtCoverage(kappa, cov)
	sat, err := coverage.Determine(cp)
	if err != nil {
		panic(err)
	}
	return Assemble(plateau, cov, gain, sat)
}

func TestCoverageDoesNotEnterBendingModel(t *testing.T) {
	// The plateau arc height is a property of the shot and the strip alone. The
	// bending model has no coverage field, which is the separation the
	// saturation rule relies on: a bug that routed coverage into the moment
	// would have to change one of the physical fields instead.
	pLow := model.BuildReport(testModelParams())
	pBis := model.BuildReport(testModelParams())
	if pLow.PlateauArcHeight != pBis.PlateauArcHeight {
		t.Errorf("plateau must be a pure function of the physical inputs")
	}
	// The assembled arc height must still rise with the coverage, and the
	// increment 0.5 -> 1.0 is smaller than the increment 0 -> 0.5.
	rHalf := fullResult(48, 1.29, 0.5, 2.6)
	rFull := fullResult(48, 1.29, 1.0, 2.6)
	firstIncrement := rHalf.ArcHeight // gain 0 -> 0.5
	secondIncrement := rFull.ArcHeight - rHalf.ArcHeight
	if !(secondIncrement > 0) {
		t.Errorf("raising coverage 0.5 -> 1.0 must raise the arc height: got %g", secondIncrement)
	}
	if !(secondIncrement < firstIncrement) {
		t.Errorf("coverage increment 0.5->1.0 (%g) must be below 0->0.5 (%g)", secondIncrement, firstIncrement)
	}
}

func TestSaturationDoesNotDependOnVelocity(t *testing.T) {
	// The saturation decision comes from the gain curve alone. Two cases with
	// different velocities but the same coverage must share the same saturation
	// state, while their arc heights differ by the pinned power law.
	low := fullResult(40, 1.29, 0.98, 2.6)
	high := fullResult(80, 1.29, 0.98, 2.6)
	if low.Saturated != high.Saturated {
		t.Errorf("saturation must not depend on velocity: low=%v high=%v", low.Saturated, high.Saturated)
	}
	if !high.Saturated {
		t.Errorf("reference coverage must be saturated at either velocity")
	}
	if !closeEnough(high.ArcHeight/low.ArcHeight, 4.0, 1e-9) {
		t.Errorf("doubling the velocity must quadruple the arc height at equal coverage: got %g", high.ArcHeight/low.ArcHeight)
	}
}

func TestSaturationDoesNotDependOnThickness(t *testing.T) {
	thin := fullResult(48, 1.29, 0.98, 2.6)
	thick := fullResult(48, 2.00, 0.98, 2.6)
	if thin.Saturated != thick.Saturated {
		t.Errorf("saturation must not depend on thickness: thin=%v thick=%v", thin.Saturated, thick.Saturated)
	}
	if !(thick.ArcHeight < thin.ArcHeight) {
		t.Errorf("thickening the strip must reduce the arc height: got %g vs %g", thick.ArcHeight, thin.ArcHeight)
	}
}

func TestUnsaturatedCaseStaysGradeFree(t *testing.T) {
	r := fullResult(48, 1.29, 0.5, 2.6)
	if r.Saturated {
		t.Errorf("coverage 0.5 must not be saturated")
	}
	if r.HasGrade() {
		t.Errorf("an unsaturated case must not carry a grade letter")
	}
	if r.GradeLetter() != "" {
		t.Errorf("grade letter must be empty for an unsaturated case, got %q", r.GradeLetter())
	}
}
