package coverage

import (
	"math"
	"testing"
)

func TestHalfLife(t *testing.T) {
	lambda := 0.085
	half := HalfLife(lambda)
	if !closeEnough(half, math.Ln2/0.085, 1e-12) {
		t.Errorf("half-life must be ln(2)/lambda: got %g", half)
	}
	if !closeEnough(CoverageAtTime(lambda, half), 0.5, 1e-12) {
		t.Errorf("coverage at the half-life must be exactly one half: got %g", CoverageAtTime(lambda, half))
	}
	if !closeEnough(CoverageAtHalfLife(), 0.5, 0) {
		t.Errorf("CoverageAtHalfLife must return 0.5")
	}
}

func TestCoverageTable(t *testing.T) {
	rows := CoverageTable(0.085, 6)
	if len(rows) != 7 {
		t.Fatalf("expected 7 rows for 6 half-lives, got %d", len(rows))
	}
	if rows[0].Coverage != 0 {
		t.Errorf("the table must start at t=0 with coverage 0")
	}
	prev := 0.0
	for _, p := range rows {
		if p.Coverage < prev {
			t.Errorf("coverage must be monotonic in the table")
		}
		prev = p.Coverage
	}
	last := rows[len(rows)-1]
	if last.Coverage < 0.98 {
		t.Errorf("six half-lives must nearly complete the coverage, got %g", last.Coverage)
	}
}

func TestTimeToCoverageValidation(t *testing.T) {
	tm, err := TimeToCoverage(0.085, 0.98)
	if err != nil {
		t.Fatalf("TimeToCoverage must accept 0.98: %v", err)
	}
	if !closeEnough(CoverageAtTime(0.085, tm), 0.98, 1e-9) {
		t.Errorf("TimeToCoverage must round-trip the coverage")
	}
	if _, err := TimeToCoverage(0, 0.5); err == nil {
		t.Errorf("a non-positive rate constant must be rejected")
	}
	if _, err := TimeToCoverage(0.085, 1.0); err == nil {
		t.Errorf("a target coverage of 1 must be rejected")
	}
	if _, err := TimeToCoverage(0.085, 0.0); err == nil {
		t.Errorf("a target coverage of 0 must be rejected")
	}
}

func TestRemainingGain(t *testing.T) {
	kappa := 2.6
	if !closeEnough(RemainingGain(kappa, 0.0), 1.0, 1e-12) {
		t.Errorf("no coverage must leave the full gain remaining")
	}
	remaining := RemainingGain(kappa, 0.98)
	gain := GainAtCoverage(kappa, 0.98)
	if !closeEnough(remaining+gain, 1.0, 1e-12) {
		t.Errorf("remaining gain and gain must sum to one")
	}
	if remaining <= 0 || remaining >= 1 {
		t.Errorf("remaining gain at 0.98 must be in (0,1), got %g", remaining)
	}
}

func TestRelativeHeadroom(t *testing.T) {
	kappa := 2.6
	if !(RelativeHeadroom(kappa, 0.5) > 0.10) {
		t.Errorf("relative headroom at 0.5 must exceed 10%%, got %g", RelativeHeadroom(kappa, 0.5))
	}
	if !(RelativeHeadroom(kappa, 0.98) < 0.10) {
		t.Errorf("relative headroom at 0.98 must sit below 10%%, got %g", RelativeHeadroom(kappa, 0.98))
	}
	if (RelativeHeadroom(kappa, 0.5) > 0.10) != !IsSaturatedAt(kappa, 0.5) {
		t.Errorf("relative headroom and saturation must agree at 0.5")
	}
	if (RelativeHeadroom(kappa, 0.98) > 0.10) != !IsSaturatedAt(kappa, 0.98) {
		t.Errorf("relative headroom and saturation must agree at 0.98")
	}
}

func TestCoverageFractionAt(t *testing.T) {
	lambda, kappa := 0.085, 2.6
	at := CoverageFractionAt(lambda, kappa, 0)
	if at != 0 {
		t.Errorf("no time must give no gain, got %g", at)
	}
	prev := 0.0
	for _, tm := range []float64{5, 20, 60, 120} {
		f := CoverageFractionAt(lambda, kappa, tm)
		if f <= prev {
			t.Errorf("coverage fraction must grow with time: got %g after %g", f, prev)
		}
		prev = f
	}
}
