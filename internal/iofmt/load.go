package iofmt

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func LoadFile(path string) (*CaseDoc, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open case file: %w", err)
	}
	defer f.Close()

	doc := &CaseDoc{}
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(doc); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return doc, nil
}

func LoadBytes(data []byte) (*CaseDoc, error) {
	doc := &CaseDoc{}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(doc); err != nil {
		return nil, fmt.Errorf("decode case document: %w", err)
	}
	return doc, nil
}

func Validate(doc *CaseDoc) []string {
	var issues []string
	if doc == nil {
		return []string{"case document is nil"}
	}
	issues = append(issues, require(doc.Shot.Velocity, fieldShotVelocity, "m/s")...)
	issues = append(issues, require(doc.Shot.Diameter, fieldShotDiameter, "mm")...)
	issues = append(issues, require(doc.Shot.Density, fieldShotDensity, "kg/m^3")...)
	issues = append(issues, require(doc.Strip.Thickness, fieldStripThickness, "mm")...)
	issues = append(issues, require(doc.Strip.Width, fieldStripWidth, "mm")...)
	issues = append(issues, require(doc.Strip.Length, fieldStripLength, "mm")...)
	issues = append(issues, require(doc.Strip.Modulus, fieldStripModulus, "GPa")...)
	issues = append(issues, require(doc.Residual.Stress, fieldResidualStress, "MPa")...)
	issues = append(issues, require(doc.Residual.LayerDepth, fieldResidualLayer, "mm")...)
	issues = append(issues, require(doc.Process.Coverage, fieldProcessCoverage, "")...)
	issues = append(issues, require(doc.Process.RateConstant, fieldProcessRate, "1/min")...)
	issues = append(issues, require(doc.Process.GainCoefficient, fieldProcessGain, "")...)
	issues = append(issues, require(doc.Reference.Velocity, fieldReferenceVelocity, "m/s")...)
	return issues
}

func require(v *float64, name, unit string) []string {
	if v != nil {
		return nil
	}
	if unit == "" {
		return []string{fmt.Sprintf("missing required field %q", name)}
	}
	return []string{fmt.Sprintf("missing required field %q (unit: %s)", name, unit)}
}

func JoinIssues(prefix string, issues []string) error {
	if len(issues) == 0 {
		return nil
	}
	var b strings.Builder
	b.WriteString(prefix)
	for i, msg := range issues {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("  - ")
		b.WriteString(msg)
	}
	return fmt.Errorf("%s", b.String())
}
