package model

const (
	mmPerMicron  = 0.001
	micronPerMM  = 1000.0
	paPerMPa     = 1e6
	mpaPerGPa    = 1000.0
	m3PerMM3     = 1e-9
	newtonPerKGF = 9.80665
)

func ToMicrometres(mm float64) float64 {
	return mm * micronPerMM
}

func ToMillimetres(um float64) float64 {
	return um * mmPerMicron
}

func ToGPa(mpa float64) float64 {
	return mpa / mpaPerGPa
}

func ToMPa(gpa float64) float64 {
	return gpa * mpaPerGPa
}

func ToPascal(mpa float64) float64 {
	return mpa * paPerMPa
}

func ToMPaFromPascal(pa float64) float64 {
	return pa / paPerMPa
}

func ModulusPerSquareMM(gpa float64) float64 {
	return gpa * mpaPerGPa
}

func CubeMillimetresToCubicMetres(mm3 float64) float64 {
	return mm3 * m3PerMM3
}

func KiloGramForceToNewton(kgf float64) float64 {
	return kgf * newtonPerKGF
}

func NewtonToKiloGramForce(n float64) float64 {
	return n / newtonPerKGF
}

func MillinewtonToNewton(mn float64) float64 {
	return mn / 1000.0
}

func NewtonMillimetreToNewtonMeter(nmm float64) float64 {
	return nmm / 1000.0
}
