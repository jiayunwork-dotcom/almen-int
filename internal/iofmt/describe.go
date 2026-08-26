package iofmt

import (
	"fmt"
	"strings"

	"almen-int/internal/model"
)

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

func DescribeParams(mp model.Params) string {
	return fmt.Sprintf("strip %.2f mm x %.1f mm x %.1f mm, E=%.0f GPa, v=%.1f m/s",
		mp.Thickness, mp.Width, mp.Length, mp.Modulus, mp.Velocity)
}

func DescribeShot(mp model.Params) string {
	return fmt.Sprintf("shot %.2f mm at %.1f m/s, density %.0f kg/m^3",
		mp.ShotDiameter, mp.Velocity, mp.ShotDensity)
}
