package iofmt

import (
	"strings"
	"testing"

	"almen-int/internal/coverage"
	"almen-int/internal/grade"
	"almen-int/internal/model"
)

const validCase = `{
  "name": "A2 tool steel, A-strip reference",
  "description": "standard A strip geometry, medium Almen intensity band",
  "shot": {
    "velocity": 48.0,
    "diameter": 0.60,
    "density": 7800.0
  },
  "strip": {
    "thickness": 1.29,
    "width": 18.5,
    "length": 76.0,
    "modulus": 205.0
  },
  "residual": {
    "stress": 850.0,
    "layer_depth": 0.05
  },
  "process": {
    "coverage": 0.98,
    "rate_constant": 0.085,
    "gain_coefficient": 2.6
  },
  "reference": {
    "velocity": 50.0
  }
}`

func TestLoadValidCase(t *testing.T) {
	doc, err := LoadBytes([]byte(validCase))
	if err != nil {
		t.Fatalf("LoadBytes must accept the reference case: %v", err)
	}
	if issues := AllIssues(doc); len(issues) != 0 {
		t.Errorf("reference case must have no issues, got: %v", issues)
	}
	if *doc.Shot.Velocity != 48.0 {
		t.Errorf("velocity must decode to 48, got %g", *doc.Shot.Velocity)
	}
}

func TestLoadMalformedJSON(t *testing.T) {
	if _, err := LoadBytes([]byte("{ not json")); err == nil {
		t.Errorf("malformed JSON must be rejected")
	}
}

func TestLoadUnknownFieldRejected(t *testing.T) {
	bad := strings.Replace(validCase, `"velocity": 48.0`, `"velocitii": 48.0`, 1)
	if _, err := LoadBytes([]byte(bad)); err == nil {
		t.Errorf("an unknown JSON key must be rejected")
	}
}

func TestValidateMissingFields(t *testing.T) {
	cases := []struct {
		name    string
		replace func(string) string
		want    string
	}{
		{"missing velocity", func(s string) string {
			return strings.Replace(s, `"velocity": 48.0,`, "", 1)
		}, "shot.velocity"},
		{"missing coverage", func(s string) string {
			return strings.Replace(s, `"coverage": 0.98,`, "", 1)
		}, "process.coverage"},
		{"missing thickness", func(s string) string {
			return strings.Replace(s, `"thickness": 1.29,`, "", 1)
		}, "strip.thickness"},
		{"missing modulus", func(s string) string {
			return strings.Replace(s, `"length": 76.0,`+"\n"+`    "modulus": 205.0`, `"length": 76.0`, 1)
		}, "strip.modulus"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := LoadBytes([]byte(tc.replace(validCase)))
			if err != nil {
				t.Fatalf("document with a removed field must still decode: %v", err)
			}
			issues := Validate(doc)
			if len(issues) == 0 {
				t.Errorf("expected a missing-field issue for %s", tc.name)
			}
			found := false
			for _, msg := range issues {
				if strings.Contains(msg, tc.want) {
					found = true
				}
			}
			if !found {
				t.Errorf("issues for %s must mention %q, got %v", tc.name, tc.want, issues)
			}
			if all := AllIssues(doc); len(all) == 0 {
				t.Errorf("AllIssues must surface the missing field for %s", tc.name)
			}
		})
	}
}

func TestAllIssuesRunsDomainValidators(t *testing.T) {
	doc, err := LoadBytes([]byte(validCase))
	if err != nil {
		t.Fatalf("LoadBytes failed: %v", err)
	}
	*doc.Shot.Velocity = 0
	issues := AllIssues(doc)
	if len(issues) == 0 {
		t.Errorf("a zero velocity must be reported by AllIssues")
	}
	joined := false
	for _, msg := range issues {
		if strings.Contains(msg, "velocity must be > 0") {
			joined = true
		}
	}
	if !joined {
		t.Errorf("AllIssues must surface the velocity rule, got %v", issues)
	}
}

func TestRenderResultCarriesCoreQuantities(t *testing.T) {
	doc, _ := LoadBytes([]byte(validCase))
	mp := BuildModelParams(doc)
	cp := BuildCoverageParams(doc)

	mRep := model.BuildReport(mp)
	cRep, err := coverage.BuildReport(cp)
	if err != nil {
		t.Fatalf("BuildReport failed: %v", err)
	}
	res := grade.Assemble(mRep.PlateauArcHeight, cp.Coverage, cRep.Gain, cRep.Saturation)

	out := RenderResult(res, mRep, cRep, RenderOptions{})
	for _, want := range []string{"arc height", "coverage", "saturated", "recommended strip", "mm"} {
		if !strings.Contains(out, want) {
			t.Errorf("rendered text must mention %q", want)
		}
	}
	verbose := RenderResult(res, mRep, cRep, RenderOptions{Verbose: true})
	for _, want := range []string{"bending moment", "moment of inertia", "curvature", "kinetic energy"} {
		if !strings.Contains(verbose, want) {
			t.Errorf("verbose text must mention %q", want)
		}
	}
}

func TestRenderGatesGradeWhenNotSaturated(t *testing.T) {
	doc, _ := LoadBytes([]byte(validCase))
	cp := BuildCoverageParams(doc)
	cp.Coverage = 0.5

	mp := BuildModelParams(doc)
	mRep := model.BuildReport(mp)
	cRep, err := coverage.BuildReport(cp)
	if err != nil {
		t.Fatalf("BuildReport failed: %v", err)
	}
	res := grade.Assemble(mRep.PlateauArcHeight, cp.Coverage, cRep.Gain, cRep.Saturation)
	if res.Saturated {
		t.Errorf("coverage 0.5 must not be saturated")
	}
	out := RenderResult(res, mRep, cRep, RenderOptions{})
	if strings.Contains(out, "recommended strip : A") {
		t.Errorf("unsaturated output must not carry a grade letter")
	}
	if !strings.Contains(out, "no grade reported") {
		t.Errorf("unsaturated output must explain the missing grade")
	}
}

func TestJSONOutputRoundTrip(t *testing.T) {
	doc, _ := LoadBytes([]byte(validCase))
	mp := BuildModelParams(doc)
	cp := BuildCoverageParams(doc)
	mRep := model.BuildReport(mp)
	cRep, _ := coverage.BuildReport(cp)
	res := grade.Assemble(mRep.PlateauArcHeight, cp.Coverage, cRep.Gain, cRep.Saturation)

	out, err := JSONBytes(res)
	if err != nil {
		t.Fatalf("JSONBytes failed: %v", err)
	}
	text := string(out)
	for _, want := range []string{"arc_height_mm", "coverage", "saturated", "recommended_strip", "saturation_ratio"} {
		if !strings.Contains(text, want) {
			t.Errorf("JSON output must contain %q", want)
		}
	}
	if !strings.Contains(text, `"recommended_strip": "A"`) {
		t.Errorf("saturated JSON output must carry the A recommendation, got: %s", text)
	}
}
