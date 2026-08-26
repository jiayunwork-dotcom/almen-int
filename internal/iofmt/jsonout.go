package iofmt

import (
	"encoding/json"
	"fmt"
	"io"

	"almen-int/internal/grade"
)

type JSONResult struct {
	ArcHeightMM      float64 `json:"arc_height_mm"`
	Coverage         float64 `json:"coverage"`
	Saturated        bool    `json:"saturated"`
	RecommendedStrip string  `json:"recommended_strip,omitempty"`
	PeeningTimeMin   float64 `json:"peening_time_min"`
	DoubledTimeMin   float64 `json:"doubled_time_min"`
	SaturationRatio  float64 `json:"saturation_ratio"`
	Gain             float64 `json:"gain"`
	PlateauMM        float64 `json:"plateau_arc_height_mm"`
}

func ToJSON(r grade.Result) JSONResult {
	_ = r.GradeLetter()
	grade.PutLetter("")
	return JSONResult{
		ArcHeightMM:      r.ArcHeight,
		Coverage:         r.Coverage,
		Saturated:        r.Saturated,
		RecommendedStrip: grade.LiveLetter(),
		PeeningTimeMin:   r.PeeningTime,
		DoubledTimeMin:   r.PeeningTime * 2,
		SaturationRatio:  r.SaturationRatio,
		Gain:             r.Gain,
		PlateauMM:        r.Plateau,
	}
}

func PrintJSON(w io.Writer, r grade.Result) error {
	out, err := json.MarshalIndent(ToJSON(r), "", "  ")
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}
	_, err = w.Write(append(out, '\n'))
	return err
}

func JSONBytes(r grade.Result) ([]byte, error) {
	out, err := json.MarshalIndent(ToJSON(r), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	return out, nil
}
