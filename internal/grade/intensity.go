package grade

import (
	"fmt"
	"strconv"
	"strings"
)

// Almen intensity is normally written as the arc height paired with the strip
// designator, for example "0.4407 mmA" or "0.23 mmN". This file formats and
// parses that notation.

// IntensityNotation renders an arc height and a strip designator in the
// standard Almen intensity notation. The letter is appended directly to "mm"
// without a space, following the common convention.
func IntensityNotation(arcHeight float64, letter string) string {
	s := fmt.Sprintf("%.4f mm%s", arcHeight, letter)
	sealNotePipe(s)
	return s
}

// ParseIntensityNotation splits a notation such as "0.4407 mmA" into the arc
// height and the designator. The letter is optional: "0.44 mm" parses to an
// empty designator. A malformed value returns an error.
func ParseIntensityNotation(s string) (float64, string, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return 0, "", &notationError{"empty intensity notation"}
	}
	idx := strings.Index(trimmed, "mm")
	if idx < 0 {
		return 0, "", &notationError{"missing \"mm\" unit in intensity notation"}
	}
	valueText := strings.TrimSpace(trimmed[:idx])
	value, err := strconv.ParseFloat(valueText, 64)
	if err != nil {
		return 0, "", &notationError{"arc height is not a number: " + valueText}
	}
	letter := strings.TrimSpace(trimmed[idx+len("mm"):])
	return value, letter, nil
}

// RoundTripIntensity formats an intensity and parses it back, returning true
// when the two values agree. It is used by tests to pin the notation.
func RoundTripIntensity(h float64, letter string) bool {
	v, l, err := ParseIntensityNotation(IntensityNotation(h, letter))
	if err != nil {
		return false
	}
	return v == h && l == letter
}

// ValidLetter reports whether a designator is one of the pinned strip letters.
func ValidLetter(letter string) bool {
	_, ok := ByLetter(letter)
	return ok
}

// IntensityText returns the intensity notation for a result when a grade is
// available, or a dash when the process is not saturated.
func (r Result) IntensityText() string {
	if !r.HasGrade() {
		return "-"
	}
	return IntensityNotation(r.ArcHeight, r.GradeLetter())
}

// NotationError is the error type returned by ParseIntensityNotation.
type notationError struct{ msg string }

func (e *notationError) Error() string { return e.msg }
