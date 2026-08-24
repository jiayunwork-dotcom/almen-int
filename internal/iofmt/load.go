package iofmt

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// LoadFile reads a case file from disk, decodes it as a CaseDoc and returns
// the document. The decoder rejects unknown fields so that a typo in a JSON
// key is reported instead of being silently ignored. The path is included in
// every error so the CLI can tell which file failed.
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

// LoadBytes decodes a case document from a byte slice. It is used by tests and
// by callers that already hold the file contents in memory.
func LoadBytes(data []byte) (*CaseDoc, error) {
	doc := &CaseDoc{}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(doc); err != nil {
		return nil, fmt.Errorf("decode case document: %w", err)
	}
	return doc, nil
}

// Validate checks the document shape: every required numeric field must be
// present. The list of issues is returned; each message names the offending
// document key.
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

// require checks one pointer field and returns a diagnostic when it is absent.
func require(v *float64, name, unit string) []string {
	if v != nil {
		return nil
	}
	if unit == "" {
		return []string{fmt.Sprintf("missing required field %q", name)}
	}
	return []string{fmt.Sprintf("missing required field %q (unit: %s)", name, unit)}
}

// JoinIssues renders a list of issues as a single error message, one issue per
// line, prefixed with the given context.
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
