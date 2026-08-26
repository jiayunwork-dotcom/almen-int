package grade

import (
	"testing"

	"almen-int/internal/coverage"
)

func testCoverageParams() coverage.Params {
	return coverage.Params{
		Coverage:        0.98,
		RateConstant:    0.085,
		GainCoefficient: 2.6,
	}
}

func TestRecommendBandSelection(t *testing.T) {
	cases := []struct {
		arc  float64
		want string
	}{
		{0.05, "N"},
		{0.10, "A"},
		{0.44, "A"},
		{0.60, "C"},
		{0.90, "C"},
	}
	for _, tc := range cases {
		r := Recommend(tc.arc, true)
		if !r.Available {
			t.Errorf("arc %g with saturation must be mappable, got unavailable", tc.arc)
		}
		if got := r.GradeLetter(); got != tc.want {
			t.Errorf("arc %g must map to strip %q, got %q", tc.arc, tc.want, got)
		}
	}
}

func TestGradeGatingWhenNotSaturated(t *testing.T) {
	r := Recommend(0.44, false)
	if r.Available {
		t.Errorf("an unsaturated process must not produce a grade letter, got %q", r.GradeLetter())
	}
	if got := r.GradeLetter(); got != "" {
		t.Errorf("grade letter must be empty when not saturated, got %q", got)
	}
	sat, err := coverage.Determine(coverage.WithExplicitCoverage(testCoverageParams(), 0.5))
	if err != nil {
		t.Fatalf("Determine failed: %v", err)
	}
	res := Assemble(0.478, 0.5, coverage.GainAtCoverage(2.6, 0.5), sat)
	if res.HasGrade() {
		t.Errorf("assembled result at low coverage must not have a grade, got %q", res.GradeLetter())
	}
}

func TestAssembleAppliesGain(t *testing.T) {
	plateau := 0.478
	gain := 0.92
	sat, err := coverage.Determine(testCoverageParams())
	if err != nil {
		t.Fatalf("Determine failed: %v", err)
	}
	res := Assemble(plateau, 0.98, gain, sat)
	want := plateau * gain
	if !closeEnough(res.ArcHeight, want, 1e-12) {
		t.Errorf("arc height must equal plateau*gain: got %g, want %g", res.ArcHeight, want)
	}
	if !res.Saturated {
		t.Errorf("reference coverage must be saturated")
	}
	if !res.HasGrade() {
		t.Errorf("saturated result must carry a grade letter")
	}
}

func TestBandsPartitionTheAxis(t *testing.T) {
	if !BandsPartition() {
		t.Errorf("strip bands must partition the non-negative axis")
	}
	if !BandsOverlapFree() {
		t.Errorf("strip bands must not overlap")
	}
	if !UniqueLetters() {
		t.Errorf("strip designators must be unique")
	}
	if len(Strips) != 3 {
		t.Errorf("expected three pinned strips, got %d", len(Strips))
	}
	a, ok := ByLetter("A")
	if !ok {
		t.Fatalf("A strip must exist")
	}
	if !closeEnough(a.Thickness, 1.29, 1e-9) {
		t.Errorf("A strip nominal thickness must be 1.29 mm, got %g", a.Thickness)
	}
}

func TestBandEdges(t *testing.T) {
	n, _ := ByLetter("N")
	a, _ := ByLetter("A")
	c, _ := ByLetter("C")
	if a.Lower != n.Upper {
		t.Errorf("A band must start where the N band ends: %g vs %g", a.Lower, n.Upper)
	}
	if c.Lower != a.Upper {
		t.Errorf("C band must start where the A band ends: %g vs %g", c.Lower, a.Upper)
	}
	if !a.BandContains(0.44) {
		t.Errorf("A band must contain 0.44 mm")
	}
	if a.BandContains(0.60) {
		t.Errorf("A band must exclude the upper edge 0.60 mm")
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
