package model

import "fmt"

// Validate checks every physical input of the bending model and returns the
// list of violations found. An empty list means the parameter set is usable.
//
// The rules follow the model contract:
//
//   - a zero or negative velocity is rejected (the model would otherwise
//     report an arc height for a process that never ran);
//   - the strip, the shot and the reference must all have positive sizes;
//   - the residual layer must be thinner than half the strip, otherwise the
//     lever-arm term (t-d) in the moment formula loses its physical meaning;
//   - the shot must be smaller than the strip so the geometry stays sensible.
//
// Callers aggregate the returned messages and surface them as a single error.
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

// Valid reports whether the parameter set passes every rule in Validate.
func Valid(p Params) bool {
	return len(Validate(p)) == 0
}

// FirstIssue returns the first violation, or an empty string when the
// parameter set is valid. It is a convenience for callers that only need a
// single diagnostic line.
func FirstIssue(p Params) string {
	issues := Validate(p)
	if len(issues) == 0 {
		return ""
	}
	return issues[0]
}
