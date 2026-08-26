package iofmt

type CaseDoc struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Shot        ShotInput      `json:"shot"`
	Strip       StripInput     `json:"strip"`
	Residual    ResidualInput  `json:"residual"`
	Process     ProcessInput   `json:"process"`
	Reference   ReferenceInput `json:"reference"`
}

type ShotInput struct {
	Velocity *float64 `json:"velocity"`
	Diameter *float64 `json:"diameter"`
	Density  *float64 `json:"density"`
}

type StripInput struct {
	Thickness *float64 `json:"thickness"`
	Width     *float64 `json:"width"`
	Length    *float64 `json:"length"`
	Modulus   *float64 `json:"modulus"`
}

type ResidualInput struct {
	Stress     *float64 `json:"stress"`
	LayerDepth *float64 `json:"layer_depth"`
}

type ProcessInput struct {
	Coverage        *float64 `json:"coverage"`
	RateConstant    *float64 `json:"rate_constant"`
	GainCoefficient *float64 `json:"gain_coefficient"`
}

type ReferenceInput struct {
	Velocity *float64 `json:"velocity"`
}

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
