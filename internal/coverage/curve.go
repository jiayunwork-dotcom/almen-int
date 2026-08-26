package coverage

import "fmt"

type IntensityPoint struct {
	Time      float64
	Coverage  float64
	Gain      float64
	Saturated bool
}

func IntensityAt(lambda, kappa, t float64) (IntensityPoint, error) {
	if lambda <= 0 {
		return IntensityPoint{}, fmt.Errorf("rate constant must be positive")
	}
	if kappa <= 0 {
		return IntensityPoint{}, fmt.Errorf("gain coefficient must be positive")
	}
	if t < 0 {
		return IntensityPoint{}, fmt.Errorf("peening time must be non-negative")
	}
	cov := CoverageAtTime(lambda, t)
	gain := GainAtCoverage(kappa, cov)
	return IntensityPoint{
		Time:      t,
		Coverage:  cov,
		Gain:      gain,
		Saturated: IsSaturatedAt(kappa, cov),
	}, nil
}

func SaturationTime(lambda, kappa float64) (float64, error) {
	if lambda <= 0 {
		return 0, fmt.Errorf("rate constant must be positive")
	}
	if kappa <= 0 {
		return 0, fmt.Errorf("gain coefficient must be positive")
	}
	cStar := ThresholdCoverage(kappa)
	return TimeForCoverage(lambda, cStar)
}

func CurveUntilSaturation(lambda, kappa float64, steps int) ([]IntensityPoint, error) {
	if steps < 2 {
		return nil, fmt.Errorf("curve needs at least two samples")
	}
	tStar, err := SaturationTime(lambda, kappa)
	if err != nil {
		return nil, err
	}
	out := make([]IntensityPoint, 0, steps)
	tEnd := tStar * (1 + 1e-6)
	for i := 0; i < steps; i++ {
		t := tEnd * float64(i) / float64(steps-1)
		pt, err := IntensityAt(lambda, kappa, t)
		if err != nil {
			return nil, err
		}
		out = append(out, pt)
	}
	return out, nil
}

func DoublingGrowthAtTime(lambda, kappa, t float64) (float64, error) {
	pt, err := IntensityAt(lambda, kappa, t)
	if err != nil {
		return 0, err
	}
	c2 := CoverageAfterDoubling(pt.Coverage)
	g2 := GainAtCoverage(kappa, c2)
	if pt.Gain == 0 {
		return 0, fmt.Errorf("gain at t=0 is zero")
	}
	return g2 / pt.Gain, nil
}

func FirstSaturatedTime(lambda, kappa float64, times []float64) (float64, bool, error) {
	if len(times) == 0 {
		return 0, false, fmt.Errorf("times must not be empty")
	}
	for _, t := range times {
		pt, err := IntensityAt(lambda, kappa, t)
		if err != nil {
			return 0, false, err
		}
		if pt.Saturated {
			return t, true, nil
		}
	}
	return 0, false, nil
}

func ArcAlongCurve(plateau, lambda, kappa float64, times []float64) ([]float64, error) {
	if plateau <= 0 {
		return nil, fmt.Errorf("plateau arc height must be positive")
	}
	out := make([]float64, 0, len(times))
	for _, t := range times {
		pt, err := IntensityAt(lambda, kappa, t)
		if err != nil {
			return nil, err
		}
		out = append(out, plateau*pt.Gain)
	}
	return out, nil
}

func SampleAroundSaturation(lambda, kappa float64) ([]IntensityPoint, error) {
	tStar, err := SaturationTime(lambda, kappa)
	if err != nil {
		return nil, err
	}
	times := []float64{0.5 * tStar, tStar * (1 + 1e-6), 2 * tStar}
	out := make([]IntensityPoint, 0, len(times))
	for _, t := range times {
		pt, err := IntensityAt(lambda, kappa, t)
		if err != nil {
			return nil, err
		}
		out = append(out, pt)
	}
	return out, nil
}
