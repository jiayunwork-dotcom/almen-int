package iofmt

import (
	"fmt"
	"strings"

	"almen-int/internal/model"
)

// DescribeCase builds a one-line human summary of a case document, listing
// the strip geometry, the shot velocity and the target coverage. It is used
// by the verbose report header so the printed numbers can be traced back to
// the input they came from.
func DescribeCase(doc *CaseDoc) string {
	var parts []string
	if doc.Name != "" {
		parts = append(parts, doc.Name)
	}
	if doc.Shot.Velocity != nil {
		parts = append(parts, fmt.Sprintf("v=%.1f m/s", *doc.Shot.Velocity))
	}
	if doc.Strip.Thickness != nil {
		parts = append(parts, fmt.Sprintf("t=%.2f mm", *doc.Strip.Thickness))
	}
	if doc.Strip.Modulus != nil {
		parts = append(parts, fmt.Sprintf("E=%.0f GPa", *doc.Strip.Modulus))
	}
	if doc.Process.Coverage != nil {
		parts = append(parts, fmt.Sprintf("coverage=%.2f", *doc.Process.Coverage))
	}
	return strings.Join(parts, ", ")
}

// DescribeParams builds a similar summary directly from the bending-model
// parameters, for use after a document has been validated and converted.
func DescribeParams(mp model.Params) string {
	return fmt.Sprintf("strip %.2f mm x %.1f mm x %.1f mm, E=%.0f GPa, v=%.1f m/s",
		mp.Thickness, mp.Width, mp.Length, mp.Modulus, mp.Velocity)
}

// DescribeShot builds a one-line summary of the shot stream.
func DescribeShot(mp model.Params) string {
	return fmt.Sprintf("shot %.2f mm at %.1f m/s, density %.0f kg/m^3",
		mp.ShotDiameter, mp.Velocity, mp.ShotDensity)
}
