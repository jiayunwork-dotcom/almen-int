package iofmt

// CaseDoc is the root document of an Almen case file. Pointers are used for
// every numeric field so that a missing field is distinguishable from a zero
// value: both fail validation, but with different messages.
type CaseDoc struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Shot        ShotInput      `json:"shot"`
	Strip       StripInput     `json:"strip"`
	Residual    ResidualInput  `json:"residual"`
	Process     ProcessInput   `json:"process"`
	Reference   ReferenceInput `json:"reference"`
}

// ShotInput describes the shot stream.
type ShotInput struct {
	Velocity *float64 `json:"velocity"` // m/s
	Diameter *float64 `json:"diameter"` // mm
	Density  *float64 `json:"density"`  // kg/m^3
}

// StripInput describes the Almen strip geometry and material.
type StripInput struct {
	Thickness *float64 `json:"thickness"` // mm
	Width     *float64 `json:"width"`     // mm
	Length    *float64 `json:"length"`    // mm
	Modulus   *float64 `json:"modulus"`   // GPa
}

// ResidualInput describes the compressive layer left by the peening.
type ResidualInput struct {
	Stress     *float64 `json:"stress"`      // MPa at the reference velocity
	LayerDepth *float64 `json:"layer_depth"` // mm
}

// ProcessInput carries the coverage target and the two curve coefficients.
type ProcessInput struct {
	Coverage        *float64 `json:"coverage"`         // in (0, 2]
	RateConstant    *float64 `json:"rate_constant"`    // 1/min
	GainCoefficient *float64 `json:"gain_coefficient"` // unitless
}

// ReferenceInput anchors the pinned velocity power law.
type ReferenceInput struct {
	Velocity *float64 `json:"velocity"` // m/s
}

// Names of the fields, reused by validation messages so that the document
// keys and the diagnostics stay in sync.
const (
	fieldShotVelocity      = "shot.velocity"
	fieldShotDiameter      = "shot.diameter"
	fieldShotDensity       = "shot.density"
	fieldStripThickness    = "strip.thickness"
	fieldStripWidth        = "strip.width"
	fieldStripLength       = "strip.length"
	fieldStripModulus      = "strip.modulus"
	fieldResidualStress    = "residual.stress"
	fieldResidualLayer     = "residual.layer_depth"
	fieldProcessCoverage   = "process.coverage"
	fieldProcessRate       = "process.rate_constant"
	fieldProcessGain       = "process.gain_coefficient"
	fieldReferenceVelocity = "reference.velocity"
)
