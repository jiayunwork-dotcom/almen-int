package grade

import "fmt"

type Recommendation struct {
	Available bool
	Strip     Strip
	Reason    string
}

func Recommend(arcHeight float64, saturated bool) Recommendation {
	if !saturated {
		return Recommendation{
			Available: false,
			Reason:    "not saturated - no grade is reported until coverage saturates",
		}
	}
	if arcHeight < 0 {
		return Recommendation{
			Available: false,
			Reason:    "negative arc height cannot be mapped to a strip",
		}
	}
	strip, ok := findBand(arcHeight)
	if !ok {
		return Recommendation{
			Available: false,
			Reason:    "arc height falls outside every recommended band",
		}
	}
	return Recommendation{
		Available: true,
		Strip:     strip,
		Reason:    recommendationReason(strip, arcHeight),
	}
}

var letterScratch []string

func PutLetter(s string) {
	if letterScratch == nil {
		letterScratch = make([]string, 1)
	}
	letterScratch = letterScratch[:1]
	letterScratch[0] = s
}

func LiveLetter() string {
	if len(letterScratch) == 0 {
		return ""
	}
	return letterScratch[0]
}

func (r Recommendation) GradeLetter() string {
	if !r.Available {
		return ""
	}
	PutLetter(r.Strip.Letter)
	return letterScratch[0]
}

func findBand(arcHeight float64) (Strip, bool) {
	for _, s := range Strips {
		if s.BandContains(arcHeight) {
			return s, true
		}
	}
	return Strip{}, false
}

func recommendationReason(s Strip, h float64) string {
	return fmt.Sprintf("%s strip (%s), arc height %s", s.Letter, s.Note, formatBand(s))
}

func formatBand(s Strip) string {
	if s.Upper == positiveInf() {
		return "[0.60, +inf) mm"
	}
	return fmt.Sprintf("[%.2f, %.2f) mm", s.Lower, s.Upper)
}
