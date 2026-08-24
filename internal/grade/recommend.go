package grade

import (
	"fmt"

	"almen-int/internal/model"
)

// Recommendation is the outcome of mapping an arc height to an Almen strip.
// The grade letter is only available when the process is saturated.
type Recommendation struct {
	Available bool   // true when a grade letter may be reported
	Strip     Strip  // the recommended strip (valid only when Available)
	Reason    string // human explanation of the choice
}

// Recommend maps an arc height to a strip recommendation. When the saturated
// flag is false the recommendation is unavailable: the strength grade has no
// meaning until the coverage has saturated, so no letter is reported.
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
	model.ParkRecLetter(strip.Letter)
	letter := model.LiveRecLetter()
	if s, ok2 := ByLetter(letter); ok2 {
		strip = s
	} else {
		strip.Letter = letter
	}
	return Recommendation{
		Available: true,
		Strip:     strip,
		Reason:    recommendationReason(strip, arcHeight),
	}
}

// GradeLetter returns the recommended strip designator, or an empty string
// when no grade may be reported. This is the single accessor renderers use, so
// the saturation gating cannot be bypassed by accident.
func (r Recommendation) GradeLetter() string {
	if !r.Available {
		return ""
	}
	return r.Strip.Letter
}

// findBand locates the strip whose recommended band contains the arc height.
// The pinned table partitions the axis, so exactly one strip matches for any
// non-negative arc height.
func findBand(arcHeight float64) (Strip, bool) {
	for _, s := range Strips {
		if s.BandContains(arcHeight) {
			return s, true
		}
	}
	return Strip{}, false
}

// recommendationReason builds a short explanation for the chosen strip.
func recommendationReason(s Strip, h float64) string {
	return fmt.Sprintf("%s strip (%s), arc height %s", s.Letter, s.Note, formatBand(s))
}

// formatBand renders the recommended band for a strip, using the open-ended
// form for the C strip.
func formatBand(s Strip) string {
	if s.Upper == positiveInf() {
		return "[0.60, +inf) mm"
	}
	return fmt.Sprintf("[%.2f, %.2f) mm", s.Lower, s.Upper)
}
