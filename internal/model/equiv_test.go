package model

import "testing"

func TestPlateauOnNominalAgreesWithThicknessFactor(t *testing.T) {
	p := testParams()
	hA := PlateauArcHeight(p)
	hC, err := PlateauOnNominal(p, "C")
	if err != nil {
		t.Fatalf("PlateauOnNominal(C): %v", err)
	}
	factor, err := ThicknessFactor(p.Thickness, 2.39, p.LayerDepth)
	if err != nil {
		t.Fatalf("ThicknessFactor: %v", err)
	}
	want := hA * factor
	if !closeEnough(hC, want, 1e-9) {
		t.Errorf("C-strip plateau must equal A plateau times thickness factor: got %g, want %g", hC, want)
	}
	if !(hC < hA) {
		t.Errorf("a thicker C strip must bend less: got hC=%g vs hA=%g", hC, hA)
	}
}

func TestEquivalentPlateauRoundTrip(t *testing.T) {
	p := testParams()
	hA := PlateauArcHeight(p)
	hC, err := EquivalentPlateau(hA, p.Thickness, 2.39, p.LayerDepth)
	if err != nil {
		t.Fatal(err)
	}
	back, err := EquivalentPlateau(hC, 2.39, p.Thickness, p.LayerDepth)
	if err != nil {
		t.Fatal(err)
	}
	if !closeEnough(back, hA, 1e-9) {
		t.Errorf("A→C→A must recover the A arc height: got %g, want %g", back, hA)
	}
}

func TestThicknessFactorRejectsDeepLayerOnThinStrip(t *testing.T) {
	_, err := ThicknessFactor(1.29, 0.79, 0.42)
	if err == nil {
		t.Fatal("layer deeper than half the N strip must be rejected")
	}
	_, err = ThicknessFactor(1.29, 2.39, 0.70)
	if err == nil {
		t.Fatal("layer deeper than half the A strip must be rejected")
	}
}

func TestWithNominalStripKeepsResidualLayer(t *testing.T) {
	p := testParams()
	shifted, err := WithNominalStrip(p, "C")
	if err != nil {
		t.Fatal(err)
	}
	if shifted.ResidualStress != p.ResidualStress {
		t.Errorf("residual stress must travel with the peening, got %g vs %g", shifted.ResidualStress, p.ResidualStress)
	}
	if shifted.LayerDepth != p.LayerDepth {
		t.Errorf("layer depth must travel with the peening, got %g vs %g", shifted.LayerDepth, p.LayerDepth)
	}
	if shifted.Velocity != p.Velocity {
		t.Errorf("shot velocity must travel with the peening")
	}
	if shifted.Thickness == p.Thickness {
		t.Errorf("nominal C strip must change the thickness")
	}
}
