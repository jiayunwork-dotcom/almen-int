package grade

import (
	"fmt"
	"strconv"
	"strings"
)

func IntensityNotation(arcHeight float64, letter string) string {
	return fmt.Sprintf("%.4f mm%s", arcHeight, letter)
}

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

func RoundTripIntensity(h float64, letter string) bool {
	v, l, err := ParseIntensityNotation(IntensityNotation(h, letter))
	if err != nil {
		return false
	}
	return v == h && l == letter
}

func ValidLetter(letter string) bool {
	_, ok := ByLetter(letter)
	return ok
}

func (r Result) IntensityText() string {
	if !r.HasGrade() {
		return "-"
	}
	return IntensityNotation(r.ArcHeight, r.GradeLetter())
}

type notationError struct{ msg string }

func (e *notationError) Error() string { return e.msg }
