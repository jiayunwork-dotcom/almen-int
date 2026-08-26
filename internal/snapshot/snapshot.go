package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"almen-int/internal/coverage"
	"almen-int/internal/grade"
	"almen-int/internal/model"
)

const (
	Magic          = "ALM1"
	CurrentVersion = 1
	Tol            = 1e-9
)

type Record struct {
	Magic             string  `json:"magic"`
	Version           int     `json:"version"`
	Velocity          float64 `json:"velocity"`
	ReferenceVelocity float64 `json:"reference_velocity"`
	ShotDiameter      float64 `json:"shot_diameter"`
	ShotDensity       float64 `json:"shot_density"`
	Thickness         float64 `json:"thickness"`
	Width             float64 `json:"width"`
	Length            float64 `json:"length"`
	Modulus           float64 `json:"modulus"`
	ResidualStress    float64 `json:"residual_stress"`
	LayerDepth        float64 `json:"layer_depth"`
	Coverage          float64 `json:"coverage"`
	RateConstant      float64 `json:"rate_constant"`
	GainCoefficient   float64 `json:"gain_coefficient"`
	ArcHeightMM       float64 `json:"arc_height_mm"`
	PlateauMM         float64 `json:"plateau_arc_height_mm"`
	Gain              float64 `json:"gain"`
	Saturated         bool    `json:"saturated"`
	GradeLetter       string  `json:"grade_letter"`
	SaturationRatio   float64 `json:"saturation_ratio"`
}

func Capture(mp model.Params, cp coverage.Params) (Record, error) {
	if issues := model.Validate(mp); len(issues) > 0 {
		return Record{}, fmt.Errorf("snapshot: invalid strip model: %s", issues[0])
	}
	cRep, err := coverage.BuildReport(cp)
	if err != nil {
		return Record{}, fmt.Errorf("snapshot: coverage: %w", err)
	}
	mRep := model.BuildReport(mp)
	res := grade.Assemble(mRep.PlateauArcHeight, cp.Coverage, cRep.Gain, cRep.Saturation)
	return fromKernel(mp, cp, res), nil
}

func fromKernel(mp model.Params, cp coverage.Params, res grade.Result) Record {
	return Record{
		Magic:             Magic,
		Version:           CurrentVersion,
		Velocity:          mp.Velocity,
		ReferenceVelocity: mp.ReferenceVelocity,
		ShotDiameter:      mp.ShotDiameter,
		ShotDensity:       mp.ShotDensity,
		Thickness:         mp.Thickness,
		Width:             mp.Width,
		Length:            mp.Length,
		Modulus:           mp.Modulus,
		ResidualStress:    mp.ResidualStress,
		LayerDepth:        mp.LayerDepth,
		Coverage:          cp.Coverage,
		RateConstant:      cp.RateConstant,
		GainCoefficient:   cp.GainCoefficient,
		ArcHeightMM:       res.ArcHeight,
		PlateauMM:         res.Plateau,
		Gain:              res.Gain,
		Saturated:         res.Saturated,
		GradeLetter:       res.GradeLetter(),
		SaturationRatio:   res.SaturationRatio,
	}
}

func (r Record) modelParams() model.Params {
	return model.Params{
		Velocity:          r.Velocity,
		ReferenceVelocity: r.ReferenceVelocity,
		ShotDiameter:      r.ShotDiameter,
		ShotDensity:       r.ShotDensity,
		Thickness:         r.Thickness,
		Width:             r.Width,
		Length:            r.Length,
		Modulus:           r.Modulus,
		ResidualStress:    r.ResidualStress,
		LayerDepth:        r.LayerDepth,
	}
}

func (r Record) coverageParams() coverage.Params {
	return coverage.Params{
		Coverage:        r.Coverage,
		RateConstant:    r.RateConstant,
		GainCoefficient: r.GainCoefficient,
	}
}

func (r Record) validate() error {
	if r.Magic != Magic {
		return fmt.Errorf("snapshot: bad magic %q", r.Magic)
	}
	if r.Version != CurrentVersion {
		return fmt.Errorf("snapshot: unsupported version %d", r.Version)
	}
	if err := r.validateInputs(); err != nil {
		return err
	}
	if r.PlateauMM <= 0 {
		return fmt.Errorf("snapshot: stored plateau must be positive")
	}
	if r.ArcHeightMM <= 0 {
		return fmt.Errorf("snapshot: stored arc height must be positive")
	}
	if r.Gain <= 0 || r.Gain >= 1 {
		return fmt.Errorf("snapshot: stored gain must lie in (0, 1)")
	}
	if r.SaturationRatio <= 0 {
		return fmt.Errorf("snapshot: stored saturation ratio must be positive")
	}
	if r.Saturated && r.GradeLetter == "" {
		return fmt.Errorf("snapshot: saturated record missing grade letter")
	}
	if !r.Saturated && r.GradeLetter != "" {
		return fmt.Errorf("snapshot: unsaturated record must not carry a grade letter")
	}
	return nil
}

func (r Record) validateInputs() error {
	mp := r.modelParams()
	if issues := model.Validate(mp); len(issues) > 0 {
		return fmt.Errorf("snapshot: %s", issues[0])
	}
	cp := r.coverageParams()
	if issues := coverage.Validate(cp); len(issues) > 0 {
		return fmt.Errorf("snapshot: %s", issues[0])
	}
	return nil
}

func WriteFile(path string, rec Record) error {
	if rec.Magic == "" {
		rec.Magic = Magic
	}
	if rec.Version == 0 {
		rec.Version = CurrentVersion
	}
	if err := rec.validate(); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	raw, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return fmt.Errorf("snapshot: empty marshal")
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func ReadFile(path string) (Record, error) {
	var rec Record
	raw, err := os.ReadFile(path)
	if err != nil {
		return rec, err
	}
	if len(raw) == 0 {
		return rec, fmt.Errorf("snapshot: empty file")
	}
	if !json.Valid(raw) {
		return rec, fmt.Errorf("snapshot: truncated or invalid JSON")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rec); err != nil {
		return Record{}, fmt.Errorf("snapshot: %w", err)
	}
	if dec.More() {
		return Record{}, fmt.Errorf("snapshot: trailing content")
	}
	if err := rec.validate(); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func (r Record) ReplayAgrees() error {
	if err := r.validate(); err != nil {
		return err
	}
	live, err := Capture(r.modelParams(), r.coverageParams())
	if err != nil {
		return err
	}
	if math.Abs(live.ArcHeightMM-r.ArcHeightMM) > Tol {
		return fmt.Errorf("snapshot: live arc height %v stored %v", live.ArcHeightMM, r.ArcHeightMM)
	}
	if math.Abs(live.PlateauMM-r.PlateauMM) > Tol {
		return fmt.Errorf("snapshot: live plateau %v stored %v", live.PlateauMM, r.PlateauMM)
	}
	if math.Abs(live.Gain-r.Gain) > Tol {
		return fmt.Errorf("snapshot: live gain %v stored %v", live.Gain, r.Gain)
	}
	if live.Saturated != r.Saturated {
		return fmt.Errorf("snapshot: live saturated %v stored %v", live.Saturated, r.Saturated)
	}
	if live.GradeLetter != r.GradeLetter {
		return fmt.Errorf("snapshot: live grade %q stored %q", live.GradeLetter, r.GradeLetter)
	}
	return nil
}

func (r Record) Matches(other Record) bool {
	if r.Magic != other.Magic || r.Version != other.Version {
		return false
	}
	if r.GradeLetter != other.GradeLetter || r.Saturated != other.Saturated {
		return false
	}
	pairs := [][2]float64{
		{r.Velocity, other.Velocity},
		{r.Thickness, other.Thickness},
		{r.LayerDepth, other.LayerDepth},
		{r.Coverage, other.Coverage},
		{r.ArcHeightMM, other.ArcHeightMM},
		{r.PlateauMM, other.PlateauMM},
		{r.Gain, other.Gain},
		{r.SaturationRatio, other.SaturationRatio},
	}
	for _, p := range pairs {
		if math.Abs(p[0]-p[1]) > Tol {
			return false
		}
	}
	return true
}
