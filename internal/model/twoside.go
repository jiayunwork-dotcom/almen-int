package model

import (
	"errors"
	"math"
)

func NetMoment(p Params, frontGain, backGain float64) (float64, error) {
	if frontGain < 0 || backGain < 0 {
		return 0, errors.New("model: coverage gain must be non-negative")
	}
	if frontGain > 1 {
		frontGain = 1
	}
	if backGain > 1 {
		backGain = 1
	}
	m := BendingMoment(p)
	net := m*frontGain - m*backGain
	if math.Abs(net) <= 1e-12*math.Abs(m) {
		return 0, nil
	}
	return net, nil
}

func TwoSidedPlateau(p Params, frontGain, backGain float64) (float64, error) {
	M, err := NetMoment(p, frontGain, backGain)
	if err != nil {
		return 0, err
	}
	den := ModulusPerArea(p) * MomentOfInertia(p)
	if den == 0 {
		return 0, errors.New("model: zero flexural rigidity")
	}
	h := (M / den) * p.Length * p.Length / 8.0
	if h < 0 {
		return -h, nil
	}
	return h, nil
}

func GainsBalanced(frontGain, backGain, tol float64) bool {
	d := frontGain - backGain
	if d < 0 {
		d = -d
	}
	return d <= tol
}

func OneSidedMatchesPlateau(p Params, frontGain float64) (float64, error) {
	h, err := TwoSidedPlateau(p, frontGain, 0)
	if err != nil {
		return 0, err
	}
	return h, nil
}

func FlipSidesNegatesMoment(p Params, frontGain, backGain float64) (bool, error) {
	a, err := NetMoment(p, frontGain, backGain)
	if err != nil {
		return false, err
	}
	b, err := NetMoment(p, backGain, frontGain)
	if err != nil {
		return false, err
	}
	return Abs(a+b) <= 1e-12, nil
}
