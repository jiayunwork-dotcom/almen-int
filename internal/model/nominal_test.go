package model

import (
	"math"
	"testing"
)

func TestNominalStripTable(t *testing.T) {
	if len(StandardStrips) != 3 {
		t.Fatalf("expected three standard strips, got %d", len(StandardStrips))
	}
	a, ok := NominalFor("A")
	if !ok {
		t.Fatalf("A strip must exist in the table")
	}
	if !closeEnough(a.Thickness, 1.29, 1e-9) {
		t.Errorf("A strip nominal thickness must be 1.29 mm, got %g", a.Thickness)
	}
	if !closeEnough(a.Width, 18.5, 1e-9) || !closeEnough(a.Length, 76.0, 1e-9) {
		t.Errorf("A strip must be 76 mm x 18.5 mm, got %g x %g", a.Length, a.Width)
	}
	if _, ok := NominalFor("Z"); ok {
		t.Errorf("an unknown designator must not resolve to a strip")
	}
}

func TestNominalStripMatches(t *testing.T) {
	p := testParams()
	// The reference case uses exactly the A-strip geometry.
	if name := StripGeometryName(p); name != "A-strip geometry" {
		t.Errorf("reference case must be recognised as A-strip geometry, got %q", name)
	}
	// A 10% thicker strip is custom geometry, not nominal.
	custom := p
	custom.Thickness = p.Thickness * 1.1
	if name := StripGeometryName(custom); name != "custom geometry" {
		t.Errorf("a non-nominal thickness must be reported as custom, got %q", name)
	}
	dev, ok := ThicknessDeviation(custom, "A")
	if !ok {
		t.Fatalf("ThicknessDeviation must resolve the A strip")
	}
	if !closeEnough(dev, p.Thickness*0.1, 1e-9) {
		t.Errorf("thickness deviation must be t - t_nominal, got %g", dev)
	}
}

func TestUnitConversions(t *testing.T) {
	if !closeEnough(ToMicrometres(1.0), 1000, 1e-12) {
		t.Errorf("1 mm must be 1000 um")
	}
	if !closeEnough(ToMillimetres(500), 0.5, 1e-12) {
		t.Errorf("500 um must be 0.5 mm")
	}
	if !closeEnough(ToGPa(205000), 205, 1e-9) {
		t.Errorf("205000 MPa must be 205 GPa")
	}
	if !closeEnough(ModulusPerSquareMM(205), 205000, 1e-9) {
		t.Errorf("205 GPa must be 205000 N/mm^2")
	}
	if !closeEnough(CubeMillimetresToCubicMetres(1), 1e-9, 1e-15) {
		t.Errorf("1 mm^3 must be 1e-9 m^3")
	}
	// The shot mass computed through the unit helper must match the direct value.
	p := testParams()
	radius := p.ShotDiameter / 2
	vol := 4.0 / 3.0 * math.Pi * radius * radius * radius
	if !closeEnough(ShotMass(p), CubeMillimetresToCubicMetres(vol)*p.ShotDensity, 1e-15) {
		t.Errorf("ShotMass must be density times the converted volume")
	}
}

func TestDeflectionRegime(t *testing.T) {
	p := testParams()
	r := EvaluateRegime(p)
	if !r.ThinStripOK {
		t.Errorf("reference case must satisfy the thin-strip check, t/L=%g", r.ThicknessOverSpan)
	}
	if !r.SmallDeflectionOK {
		t.Errorf("reference case must satisfy the small-deflection check, h/L=%g", r.ArcHeightOverSpan)
	}
	if r.SurfaceStrain <= 0 {
		t.Errorf("surface strain must be positive, got %g", r.SurfaceStrain)
	}
	if got := r.RegimeSummary(); got != "within the thin-strip and small-deflection regime" {
		t.Errorf("reference case must be inside the model regime, got %q", got)
	}
	// A very thick strip fails the thin-strip check.
	thick := p
	thick.Thickness = 10.0
	thick.Length = 76.0
	reg := EvaluateRegime(thick)
	if reg.ThinStripOK {
		t.Errorf("a 10 mm strip over a 76 mm span must fail the thin-strip check")
	}
}

func TestRegimePeakBendingStress(t *testing.T) {
	p := testParams()
	r := EvaluateRegime(p)
	// Peak bending stress is E*t*kappa/2; recompute it from the curvature.
	want := ModulusPerArea(p) * p.Thickness * Curvature(p) / 2
	if !closeEnough(r.PeakBendingStress, want, 1e-9) {
		t.Errorf("peak bending stress must equal E*t*kappa/2: got %g, want %g", r.PeakBendingStress, want)
	}
	if !(r.ElasticOK == (r.PeakBendingStress < p.ResidualStress)) {
		t.Errorf("elastic check must compare the peak stress against the residual stress")
	}
}
