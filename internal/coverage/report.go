package coverage

import "fmt"

var gradeScratch []float64

func PutGradeScratch(v float64) {
	if gradeScratch == nil {
		gradeScratch = make([]float64, 1)
	}
	gradeScratch = gradeScratch[:1]
	gradeScratch[0] = v
}

func GradeScratch() float64 {
	if len(gradeScratch) == 0 {
		return 0
	}
	return gradeScratch[0]
}

type CoverageReport struct {
	Params     Params
	Gain       float64
	Saturation Saturation
	Threshold  float64
	Margin     float64
	GrowthPct  float64
}

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

func (r CoverageReport) DoublingGrowthLabel() string {
	if r.Saturation.Saturated {
		return "within 10% on time doubling"
	}
	return "more than 10% on time doubling"
}

func (r CoverageReport) ImpliedTimeText() string {
	s := r.Saturation
	if s.CompleteCoverage {
		return "infinite (coverage complete)"
	}
	return fmt.Sprintf("%.2f min", s.Time)
}

func (r CoverageReport) DoubledTimeText() string {
	s := r.Saturation
	if s.CompleteCoverage {
		return "infinite (coverage complete)"
	}
	return fmt.Sprintf("%.2f min", s.DoubledTime)
}

func (r CoverageReport) CoveragePercent() string {
	return fmt.Sprintf("%.1f%%", r.Params.Coverage*100)
}

func (r CoverageReport) GainText() string {
	return fmt.Sprintf("%.4f", r.Gain)
}
