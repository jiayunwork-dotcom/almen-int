package coverage

import "fmt"

// CoverageReport bundles the process inputs, the gain they produce, and the
// saturation outcome into one snapshot for rendering.
type CoverageReport struct {
	Params     Params
	Gain       float64
	Saturation Saturation
	Threshold  float64
	Margin     float64
	GrowthPct  float64
}

// BuildReport runs the gain and saturation logic for the process inputs and
// returns the snapshot. The parameters are validated first; a violation is
// returned as an error.
func BuildReport(p Params) (CoverageReport, error) {
	sat, err := Determine(p)
	if err != nil {
		return CoverageReport{}, err
	}
	return CoverageReport{
		Params:     p,
		Gain:       sat.Gain,
		Saturation: sat,
		Threshold:  sat.Threshold,
		Margin:     sat.Margin(),
		GrowthPct:  sat.GrowthPercent(),
	}, nil
}

// DoublingGrowthLabel returns a short human label for the saturation state.
// The strings are used by the text renderer.
func (r CoverageReport) DoublingGrowthLabel() string {
	if r.Saturation.Saturated {
		return "within 10% on time doubling"
	}
	return "more than 10% on time doubling"
}

// ImpliedTimeText returns a printable form of the implied peening time,
// handling the infinite value that a complete coverage produces.
func (r CoverageReport) ImpliedTimeText() string {
	s := r.Saturation
	if s.CompleteCoverage {
		return "infinite (coverage complete)"
	}
	return fmt.Sprintf("%.2f min", s.Time)
}

// DoubledTimeText returns the printable doubled peening time.
func (r CoverageReport) DoubledTimeText() string {
	s := r.Saturation
	if s.CompleteCoverage {
		return "infinite (coverage complete)"
	}
	return fmt.Sprintf("%.2f min", s.DoubledTime)
}

// CoveragePercent renders the coverage as a percentage with one decimal, e.g.
// 98.0% for a coverage of 0.98.
func (r CoverageReport) CoveragePercent() string {
	return fmt.Sprintf("%.1f%%", r.Params.Coverage*100)
}

// GainText renders the arc-height gain with four decimals.
func (r CoverageReport) GainText() string {
	return fmt.Sprintf("%.4f", r.Gain)
}
