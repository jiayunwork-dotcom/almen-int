package model

import "testing"

// testParams returns a parameter set in the A-strip magnitude band. The values
// mirror example/a2-steel.json so the unit-level tests and the case file stay
// consistent with each other.
func testParams() Params {
	return Params{
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

func TestVelocityPowerLaw(t *testing.T) {
	p := testParams()
	slow := p
	slow.Velocity = 40
	fast := p
	fast.Velocity = 80

	hSlow := PlateauArcHeight(slow)
	hFast := PlateauArcHeight(fast)
	got := hFast / hSlow
	want := 4.0
	if !closeEnough(got, want, 1e-9) {
		t.Errorf("doubling the velocity must quadruple the plateau arc height: got %g, want %g", got, want)
	}

	keSlow := KineticEnergy(slow)
	keFast := KineticEnergy(fast)
	keRatio := keFast / keSlow
	if !closeEnough(keRatio, 4.0, 1e-9) {
		t.Errorf("doubling the velocity must quadruple the kinetic energy: got %g, want %g", keRatio, 4.0)
	}
}

func TestVelocityIncreaseRaisesArcHeight(t *testing.T) {
	p := testParams()
	low := p
	low.Velocity = 42
	high := p
	high.Velocity = 60

	hLow := PlateauArcHeight(low)
	hHigh := PlateauArcHeight(high)
	if !(hHigh > hLow) {
		t.Errorf("higher velocity must raise the arc height: got h=%g at v=%g vs h=%g at v=%g", hHigh, high.Velocity, hLow, low.Velocity)
	}
	if hHigh <= 0 || hLow <= 0 {
		t.Errorf("arc heights must be positive, got %g and %g", hLow, hHigh)
	}
}

func TestThickerStripReducesArcHeight(t *testing.T) {
	p := testParams()
	thin := p
	thin.Thickness = 1.29
	thick := p
	thick.Thickness = 2.0

	hThin := PlateauArcHeight(thin)
	hThick := PlateauArcHeight(thick)
	if !(hThick < hThin) {
		t.Errorf("a thicker strip must bend less: got h=%g at t=%g vs h=%g at t=%g", hThick, thick.Thickness, hThin, thin.Thickness)
	}
	if !(MomentOfInertia(thick) > MomentOfInertia(thin)) {
		t.Errorf("the inertia must grow with the thickness: got I=%g vs I=%g", MomentOfInertia(thick), MomentOfInertia(thin))
	}
}

func TestMomentOfInertiaIsCubic(t *testing.T) {
	got := InertiaScaling(1.0, 2.0)
	if !closeEnough(got, 8.0, 1e-12) {
		t.Errorf("doubling the thickness must raise the inertia by eight: got %g, want 8", got)
	}
	want := MomentOfInertiaFor(18.5, 1.29)
	direct := 18.5 * (1.29 * 1.29 * 1.29) / 12.0
	if !closeEnough(want, direct, 1e-12) {
		t.Errorf("inertia must equal w*t^3/12: got %g, want %g", want, direct)
	}
}

func TestBendingMomentFormula(t *testing.T) {
	p := testParams()
	sigma := ResidualStress(p)
	want := sigma * p.Width * p.LayerDepth * (p.Thickness - p.LayerDepth) / 2.0
	got := BendingMoment(p)
	if !closeEnough(got, want, 1e-12) {
		t.Errorf("bending moment must equal sigma*w*d*(t-d)/2: got %g, want %g", got, want)
	}
	// The moment must scale linearly with the stress, independent of the geometry.
	atStress := BendingMomentAtStress(p, sigma*2)
	if !closeEnough(atStress, got*2, 1e-9) {
		t.Errorf("doubling the stress must double the moment: got %g, want %g", atStress, got*2)
	}
}

func TestPlateauArcHeightScalesWithMoment(t *testing.T) {
	p := testParams()
	base := PlateauArcHeight(p)
	stronger := p
	stronger.ResidualStress = p.ResidualStress * 3
	if !closeEnough(PlateauArcHeight(stronger), base*3, 1e-9) {
		t.Errorf("arc height must scale linearly with the moment: got %g, want %g", PlateauArcHeight(stronger), base*3)
	}
	stiffer := p
	stiffer.Modulus = p.Modulus * 2
	if !closeEnough(PlateauArcHeight(stiffer), base/2, 1e-9) {
		t.Errorf("arc height must be inversely proportional to the modulus: got %g, want %g", PlateauArcHeight(stiffer), base/2)
	}
}

func TestModelValidationRejectsInvalidInputs(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Params)
	}{
		{"zero velocity", func(p *Params) { p.Velocity = 0 }},
		{"negative velocity", func(p *Params) { p.Velocity = -10 }},
		{"zero thickness", func(p *Params) { p.Thickness = 0 }},
		{"negative thickness", func(p *Params) { p.Thickness = -0.5 }},
		{"zero modulus", func(p *Params) { p.Modulus = 0 }},
		{"negative modulus", func(p *Params) { p.Modulus = -205 }},
		{"zero coverage stress", func(p *Params) { p.ResidualStress = 0 }},
		{"layer deeper than half strip", func(p *Params) { p.LayerDepth = p.Thickness * 0.6 }},
		{"zero layer depth", func(p *Params) { p.LayerDepth = 0 }},
		{"shot larger than strip", func(p *Params) { p.ShotDiameter = p.Thickness * 2 }},
		{"zero width", func(p *Params) { p.Width = 0 }},
		{"zero length", func(p *Params) { p.Length = 0 }},
		{"zero reference velocity", func(p *Params) { p.ReferenceVelocity = 0 }},
		{"zero density", func(p *Params) { p.ShotDensity = 0 }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := testParams()
			tc.mutate(&p)
			issues := Validate(p)
			if len(issues) == 0 {
				t.Errorf("expected validation issues for %s, got none", tc.name)
			}
			if Valid(p) {
				t.Errorf("Valid() must be false for %s", tc.name)
			}
		})
	}
}

func TestModelValidationAcceptsValidInputs(t *testing.T) {
	p := testParams()
	if issues := Validate(p); len(issues) != 0 {
		t.Errorf("reference parameter set must be valid, got issues: %v", issues)
	}
	if !Valid(p) {
		t.Errorf("Valid() must be true for the reference set")
	}
	// The maximum allowed coverage-related geometry: layer just below half the strip.
	edge := p
	edge.LayerDepth = p.Thickness/2 - 1e-6
	if issues := Validate(edge); len(issues) != 0 {
		t.Errorf("a layer just below half the thickness must be valid, got: %v", issues)
	}
}

func TestPlateauArcHeightInABand(t *testing.T) {
	p := testParams()
	h := PlateauArcHeight(p)
	if h < 0.10 || h > 0.60 {
		t.Errorf("A-strip reference must land in the 0.1-0.6 mm band, got %g mm", h)
	}
	if RelativeError(p) > 2e-3 {
		t.Errorf("small-deflection sagitta must be close to the exact value, relative error %g", RelativeError(p))
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
