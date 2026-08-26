package model

import "fmt"

func Validate(p Params) []string {
	var issues []string
	if p.Velocity <= 0 {
		issues = append(issues, fmt.Sprintf("shot velocity must be > 0 (got %g m/s)", p.Velocity))
	}
	if p.ReferenceVelocity <= 0 {
		issues = append(issues, fmt.Sprintf("reference velocity must be > 0 (got %g m/s)", p.ReferenceVelocity))
	}
	if p.ShotDiameter <= 0 {
		issues = append(issues, fmt.Sprintf("shot diameter must be > 0 (got %g mm)", p.ShotDiameter))
	}
	if p.ShotDensity <= 0 {
		issues = append(issues, fmt.Sprintf("shot density must be > 0 (got %g kg/m^3)", p.ShotDensity))
	}
	if p.Thickness <= 0 {
		issues = append(issues, fmt.Sprintf("strip thickness must be > 0 (got %g mm)", p.Thickness))
	}
	if p.Width <= 0 {
		issues = append(issues, fmt.Sprintf("strip width must be > 0 (got %g mm)", p.Width))
	}
	if p.Length <= 0 {
		issues = append(issues, fmt.Sprintf("strip length must be > 0 (got %g mm)", p.Length))
	}
	if p.Modulus <= 0 {
		issues = append(issues, fmt.Sprintf("elastic modulus must be > 0 (got %g GPa)", p.Modulus))
	}
	if p.ResidualStress <= 0 {
		issues = append(issues, fmt.Sprintf("residual stress must be > 0 (got %g MPa)", p.ResidualStress))
	}
	if p.LayerDepth <= 0 {
		issues = append(issues, fmt.Sprintf("residual layer depth must be > 0 (got %g mm)", p.LayerDepth))
	}
	if p.LayerDepth > 0 && p.Thickness > 0 && p.LayerDepth >= p.Thickness/2 {
		issues = append(issues, fmt.Sprintf(
			"residual layer depth %.3f mm must be below half the strip thickness %.3f mm",
			p.LayerDepth, p.Thickness/2,
		))
	}
	if p.ShotDiameter > 0 && p.Thickness > 0 && p.ShotDiameter >= p.Thickness {
		issues = append(issues, fmt.Sprintf(
			"shot diameter %.3f mm must be smaller than the strip thickness %.3f mm",
			p.ShotDiameter, p.Thickness,
		))
	}
	return issues
}

func Valid(p Params) bool {
	return len(Validate(p)) == 0
}

func FirstIssue(p Params) string {
	issues := Validate(p)
	if len(issues) == 0 {
		return ""
	}
	return issues[0]
}
