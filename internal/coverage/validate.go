package coverage

import "fmt"

// Validate checks the process inputs and returns the list of violations. The
// rules implement the contract that a coverage outside (0, 2] is an error and
// that both curve coefficients must be positive.
func Validate(p Params) []string {
	var issues []string
	if p.Coverage <= 0 || p.Coverage > 2 {
		issues = append(issues, fmt.Sprintf("coverage must be in (0, 2] (got %g)", p.Coverage))
	}
	if p.RateConstant <= 0 {
		issues = append(issues, fmt.Sprintf("rate constant must be > 0 (got %g 1/min)", p.RateConstant))
	}
	if p.GainCoefficient <= 0 {
		issues = append(issues, fmt.Sprintf("gain coefficient must be > 0 (got %g)", p.GainCoefficient))
	}
	return issues
}

// Valid reports whether the process inputs pass every rule in Validate.
func Valid(p Params) bool {
	return len(Validate(p)) == 0
}

// FirstIssue returns the first violation, or an empty string when the inputs
// are valid.
func FirstIssue(p Params) string {
	issues := Validate(p)
	if len(issues) == 0 {
		return ""
	}
	return issues[0]
}
