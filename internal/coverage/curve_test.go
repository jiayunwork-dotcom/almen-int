package coverage

import (
	"math"
	"testing"
)

func TestSaturationTimeAgreesWithThreshold(t *testing.T) {
	lambda := 0.085
	kappa := 2.6
	tStar, err := SaturationTime(lambda, kappa)
	if err != nil {
		t.Fatal(err)
	}
	cStar := ThresholdCoverage(kappa)
	got := CoverageAtTime(lambda, tStar)
	if !closeEnough(got, cStar, 1e-9) {
		t.Errorf("saturation time must land on the threshold coverage: got %g, want %g", got, cStar)
	}
	below, err := IntensityAt(lambda, kappa, tStar*0.98)
	if err != nil {
		t.Fatal(err)
	}
	if below.Saturated {
		t.Errorf("just below saturation time must not be saturated (coverage %g)", below.Coverage)
	}
	above, err := IntensityAt(lambda, kappa, tStar*1.02)
	if err != nil {
		t.Fatal(err)
	}
	if !above.Saturated {
		t.Errorf("just above saturation time must be saturated (coverage %g)", above.Coverage)
	}
}

func TestCurveUntilSaturationMonotone(t *testing.T) {
	pts, err := CurveUntilSaturation(0.085, 2.6, 8)
	if err != nil {
		t.Fatal(err)
	}
	if len(pts) != 8 {
		t.Fatalf("expected 8 samples, got %d", len(pts))
	}
	if pts[0].Coverage != 0 {
		t.Errorf("curve must start at coverage 0, got %g", pts[0].Coverage)
	}
	for i := 1; i < len(pts); i++ {
		if !(pts[i].Coverage > pts[i-1].Coverage) {
			t.Errorf("coverage must rise: %g then %g", pts[i-1].Coverage, pts[i].Coverage)
		}
		if !(pts[i].Gain > pts[i-1].Gain) {
			t.Errorf("gain must rise: %g then %g", pts[i-1].Gain, pts[i].Gain)
		}
	}
	last := pts[len(pts)-1]
	if !last.Saturated {
		t.Errorf("last sample must sit at saturation")
	}
}

func TestSampleAroundSaturationBracketsThreshold(t *testing.T) {
	pts, err := SampleAroundSaturation(0.085, 2.6)
	if err != nil {
		t.Fatal(err)
	}
	if len(pts) != 3 {
		t.Fatalf("expected three samples, got %d", len(pts))
	}
	if pts[0].Saturated {
		t.Errorf("half saturation time must not be saturated")
	}
	if !pts[1].Saturated {
		t.Errorf("saturation time sample must be saturated")
	}
	if !pts[2].Saturated {
		t.Errorf("double saturation time must remain saturated")
	}
	g, err := DoublingGrowthAtTime(0.085, 2.6, pts[1].Time)
	if err != nil {
		t.Fatal(err)
	}
	if g > SaturationRatio+1e-6 {
		t.Errorf("growth at saturation time must sit at the 10%% rule, got %g", g)
	}
}

func TestArcAlongCurveScalesPlateau(t *testing.T) {
	plateau := 0.478
	times := []float64{0, 10, 40}
	arcs, err := ArcAlongCurve(plateau, 0.085, 2.6, times)
	if err != nil {
		t.Fatal(err)
	}
	if len(arcs) != 3 {
		t.Fatalf("expected 3 heights, got %d", len(arcs))
	}
	if arcs[0] != 0 {
		t.Errorf("arc height at t=0 must be 0, got %g", arcs[0])
	}
	if !(arcs[1] > arcs[0] && arcs[2] > arcs[1]) {
		t.Errorf("arc height must rise along the curve: %v", arcs)
	}
	if arcs[2] >= plateau {
		t.Errorf("finite-time height must stay below the plateau, got %g vs %g", arcs[2], plateau)
	}
}

func TestFirstSaturatedTime(t *testing.T) {
	tStar, err := SaturationTime(0.085, 2.6)
	if err != nil {
		t.Fatal(err)
	}
	got, ok, err := FirstSaturatedTime(0.085, 2.6, []float64{tStar * 0.5, tStar * 0.9, tStar * 1.2})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("a sample past saturation time must be found")
	}
	if math.Abs(got-tStar*1.2) > 1e-12 {
		t.Errorf("first saturated sample must be 1.2 T*, got %g", got)
	}
	_, ok, err = FirstSaturatedTime(0.085, 2.6, []float64{tStar * 0.2, tStar * 0.4})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("samples entirely below saturation must not report a time")
	}
}

func TestIntensityAtRejectsBadInput(t *testing.T) {
	if _, err := IntensityAt(0, 2.6, 1); err == nil {
		t.Fatal("zero rate constant must fail")
	}
	if _, err := IntensityAt(0.085, 0, 1); err == nil {
		t.Fatal("zero gain coefficient must fail")
	}
	if _, err := IntensityAt(0.085, 2.6, -1); err == nil {
		t.Fatal("negative time must fail")
	}
}
