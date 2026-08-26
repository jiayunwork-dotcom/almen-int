package grade

type Strip struct {
	Letter    string
	Thickness float64
	Lower     float64
	Upper     float64
	Note      string
}

var Strips = []Strip{
	{
		Letter:    "N",
		Thickness: 0.79,
		Lower:     0.00,
		Upper:     0.10,
		Note:      "thin strip, low intensity",
	},
	{
		Letter:    "A",
		Thickness: 1.29,
		Lower:     0.10,
		Upper:     0.60,
		Note:      "standard strip, medium intensity",
	},
	{
		Letter:    "C",
		Thickness: 2.39,
		Lower:     0.60,
		Upper:     positiveInf(),
		Note:      "thick strip, high intensity",
	},
}

func ByLetter(letter string) (Strip, bool) {
	for _, s := range Strips {
		if s.Letter == letter {
			return s, true
		}
	}
	return Strip{}, false
}

func (s Strip) BandContains(h float64) bool {
	return h >= s.Lower && h < s.Upper
}

func (s Strip) BandWidth() float64 {
	return s.Upper - s.Lower
}

func (s Strip) BandCenter() float64 {
	if s.Upper == positiveInf() {
		return s.Lower
	}
	return (s.Lower + s.Upper) / 2
}

func Overlaps(a, b Strip) bool {
	return a.Lower < b.Upper && b.Lower < a.Upper
}

func positiveInf() float64 {
	return 1e300
}
