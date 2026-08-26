package grade

import (
	"testing"

	"almen-int/internal/coverage"
	"almen-int/internal/model"
)

func TestEquivalentHeightMatchesBendingModel(t *testing.T) {
	p := testModelParams()
	converted, direct, err := ConvertMatchesModel(p, "C")
	if err != nil {
		t.Fatalf("ConvertMatchesModel: %v", err)
	}
	if !closeEnough(converted, direct, 1e-9) {
		t.Errorf("letter conversion must agree with evaluating the same residual layer on a C strip: converted %g, model %g", converted, direct)
	}
	if !(converted < model.PlateauArcHeight(p)) {
		t.Errorf("C-strip equivalent height must be below the A-strip reading")
	}
}

func TestConvertAToCAndBack(t *testing.T) {
	hA := 0.4407
	hC, err := ConvertLetter(hA, "A", "C", 0.05)
	if err != nil {
		t.Fatal(err)
	}
	back, err := ConvertLetter(hC, "C", "A", 0.05)
	if err != nil {
		t.Fatal(err)
	}
	if !closeEnough(back, hA, 1e-9) {
		t.Errorf("A→C→A must recover the A reading: got %g, want %g", back, hA)
	}
}

func TestConversionCommutesWithCoverageGain(t *testing.T) {
	p := testModelParams()
	kappa := 2.6
	plateau := model.PlateauArcHeight(p)
	gain := coverage.GainAtCoverage(kappa, 0.5)
	gained := plateau * gain
	convPlateau, err := ConvertLetter(plateau, "A", "C", p.LayerDepth)
	if err != nil {
		t.Fatal(err)
	}
	convGained, err := ConvertLetter(gained, "A", "C", p.LayerDepth)
	if err != nil {
		t.Fatal(err)
	}
	if !closeEnough(convGained, convPlateau*gain, 1e-9) {
		t.Errorf("converting a gained height must equal gain times converted plateau: got %g, want %g", convGained, convPlateau*gain)
	}
}

func TestConvertedRecommendationUsesTargetBand(t *testing.T) {
	p := testModelParams()
	hA := model.PlateauArcHeight(p)
	rec, err := ConvertedRecommendation(hA, "A", "C", p.LayerDepth, true)
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Available {
		t.Fatalf("saturated conversion must remain mappable")
	}
	hC, err := ConvertLetter(hA, "A", "C", p.LayerDepth)
	if err != nil {
		t.Fatal(err)
	}
	direct := Recommend(hC, true)
	if rec.GradeLetter() != direct.GradeLetter() {
		t.Errorf("converted recommendation must use the target-strip band, got %q vs %q", rec.GradeLetter(), direct.GradeLetter())
	}
}

func TestConvertUnknownLetter(t *testing.T) {
	if _, err := ConvertLetter(0.4, "Q", "A", 0.05); err == nil {
		t.Fatal("unknown source letter must fail")
	}
	if _, err := ConvertLetter(0.4, "A", "Z", 0.05); err == nil {
		t.Fatal("unknown target letter must fail")
	}
}
