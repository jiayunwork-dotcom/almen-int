package model

import "fmt"

func WithNominalStrip(p Params, letter string) (Params, error) {
	nom, ok := NominalFor(letter)
	if !ok {
		return Params{}, fmt.Errorf("unknown strip letter %q", letter)
	}
	out := p
	out.Thickness = nom.Thickness
	out.Width = nom.Width
	out.Length = nom.Length
	if issues := Validate(out); len(issues) > 0 {
		return Params{}, fmt.Errorf("nominal %s strip is incompatible with this residual layer: %s", letter, issues[0])
	}
	return out, nil
}

func PlateauOnNominal(p Params, letter string) (float64, error) {
	shifted, err := WithNominalStrip(p, letter)
	if err != nil {
		return 0, err
	}
	return PlateauArcHeight(shifted), nil
}

func ThicknessFactor(fromThickness, toThickness, layerDepth float64) (float64, error) {
	if fromThickness <= 0 || toThickness <= 0 {
		return 0, fmt.Errorf("strip thickness must be positive")
	}
	if layerDepth <= 0 {
		return 0, fmt.Errorf("layer depth must be positive")
	}
	if layerDepth >= fromThickness/2 {
		return 0, fmt.Errorf("layer depth %g mm must stay below half the source thickness %g mm", layerDepth, fromThickness)
	}
	if layerDepth >= toThickness/2 {
		return 0, fmt.Errorf("layer depth %g mm must stay below half the target thickness %g mm", layerDepth, toThickness)
	}
	numer := (toThickness - layerDepth) / (toThickness * toThickness * toThickness)
	denom := (fromThickness - layerDepth) / (fromThickness * fromThickness * fromThickness)
	if denom == 0 {
		return 0, fmt.Errorf("source thickness factor is zero")
	}
	return numer / denom, nil
}

func EquivalentPlateau(h, fromThickness, toThickness, layerDepth float64) (float64, error) {
	if h < 0 {
		return 0, fmt.Errorf("arc height must be non-negative")
	}
	factor, err := ThicknessFactor(fromThickness, toThickness, layerDepth)
	if err != nil {
		return 0, err
	}
	return h * factor, nil
}

func SameLayerOn(p Params, letter string) (Params, float64, error) {
	shifted, err := WithNominalStrip(p, letter)
	if err != nil {
		return Params{}, 0, err
	}
	factor, err := ThicknessFactor(p.Thickness, shifted.Thickness, p.LayerDepth)
	if err != nil {
		return Params{}, 0, err
	}
	return shifted, factor, nil
}
