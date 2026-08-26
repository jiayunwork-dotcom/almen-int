package grade

import (
	"fmt"

	"almen-int/internal/model"
)

func ConvertLetter(h float64, fromLetter, toLetter string, layerDepth float64) (float64, error) {
	from, ok := ByLetter(fromLetter)
	if !ok {
		return 0, fmt.Errorf("unknown source strip %q", fromLetter)
	}
	to, ok := ByLetter(toLetter)
	if !ok {
		return 0, fmt.Errorf("unknown target strip %q", toLetter)
	}
	return model.EquivalentPlateau(h, from.Thickness, to.Thickness, layerDepth)
}

func ConvertMatchesModel(p model.Params, toLetter string) (float64, float64, error) {
	letter := stripLetterOf(p)
	if letter == "" {
		return 0, 0, fmt.Errorf("source geometry is not a standard N/A/C strip")
	}
	hFrom := model.PlateauArcHeight(p)
	converted, err := ConvertLetter(hFrom, letter, toLetter, p.LayerDepth)
	if err != nil {
		return 0, 0, err
	}
	direct, err := model.PlateauOnNominal(p, toLetter)
	if err != nil {
		return 0, 0, err
	}
	return converted, direct, nil
}

func stripLetterOf(p model.Params) string {
	for _, s := range model.StandardStrips {
		if s.Matches(p) {
			return s.Letter
		}
	}
	return ""
}

func ConvertedRecommendation(h float64, fromLetter, toLetter string, layerDepth float64, saturated bool) (Recommendation, error) {
	eq, err := ConvertLetter(h, fromLetter, toLetter, layerDepth)
	if err != nil {
		return Recommendation{}, err
	}
	return Recommend(eq, saturated), nil
}
