package model

// This file holds the unit conversions that keep the model's mixed unit
// conventions explicit. The bending formula works in N, mm and N/mm^2; the
// case files and the reports speak in GPa, MPa, mm and metres per second.
// Keeping every conversion in one file makes it impossible for the report
// layer and the tests to disagree about a factor.

// millimetre to micrometre and back.
const (
	mmPerMicron  = 0.001
	micronPerMM  = 1000.0
	paPerMPa     = 1e6
	mpaPerGPa    = 1000.0
	m3PerMM3     = 1e-9
	newtonPerKGF = 9.80665
)

// ToMicrometres converts a millimetre length to micrometres.
func ToMicrometres(mm float64) float64 {
	return mm * micronPerMM
}

// ToMillimetres converts a micrometre length to millimetres.
func ToMillimetres(um float64) float64 {
	return um * mmPerMicron
}

// ToGPa converts a stress in megapascals to gigapascals.
func ToGPa(mpa float64) float64 {
	return mpa / mpaPerGPa
}

// ToMPa converts a stress in gigapascals to megapascals.
func ToMPa(gpa float64) float64 {
	return gpa * mpaPerGPa
}

// ToPascal converts a stress in megapascals to pascals.
func ToPascal(mpa float64) float64 {
	return mpa * paPerMPa
}

// ToMPaFromPascal converts a stress in pascals to megapascals.
func ToMPaFromPascal(pa float64) float64 {
	return pa / paPerMPa
}

// ModulusPerSquareMM converts an elastic modulus in gigapascals to N/mm^2,
// the stress-like unit the bending formula needs.
func ModulusPerSquareMM(gpa float64) float64 {
	return gpa * mpaPerGPa
}

// CubeMillimetresToCubicMetres converts a volume from mm^3 to m^3. The shot
// mass is computed in this unit because the density is quoted in kg/m^3.
func CubeMillimetresToCubicMetres(mm3 float64) float64 {
	return mm3 * m3PerMM3
}

// KiloGramForceToNewton converts a force expressed in kilogram-force to
// newtons. It is provided for reports that compare the layer force against
// practical values quoted in kgf.
func KiloGramForceToNewton(kgf float64) float64 {
	return kgf * newtonPerKGF
}

// NewtonToKiloGramForce converts a force in newtons to kilogram-force.
func NewtonToKiloGramForce(n float64) float64 {
	return n / newtonPerKGF
}

// MillinewtonToNewton converts a moment-adjacent force given in millinewtons
// to newtons.
func MillinewtonToNewton(mn float64) float64 {
	return mn / 1000.0
}

// NewtonMillimetreToNewtonMeter converts a bending moment from N*mm to N*m.
func NewtonMillimetreToNewtonMeter(nmm float64) float64 {
	return nmm / 1000.0
}
