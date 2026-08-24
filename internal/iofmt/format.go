package iofmt

import "fmt"

// This file centralises the number formatting used by the text renderer, so
// every quantity is printed with the same precision everywhere.

// F4 formats a value with four decimals, the precision used for lengths.
func F4(v float64) string { return fmt.Sprintf("%.4f", v) }

// F2 formats a value with two decimals, the precision used for stresses and
// times.
func F2(v float64) string { return fmt.Sprintf("%.2f", v) }

// F6 formats a value with six decimals, the precision used for inertia and
// gain.
func F6(v float64) string { return fmt.Sprintf("%.6f", v) }

// FG formats a value with six significant digits, used for tiny quantities
// like the curvature.
func FG(v float64) string { return fmt.Sprintf("%.6g", v) }

// Pct1 formats a fraction as a percentage with one decimal.
func Pct1(frac float64) string { return fmt.Sprintf("%.1f%%", frac*100) }

// Pct2 formats a fraction as a percentage with two decimals.
func Pct2(frac float64) string { return fmt.Sprintf("%.2f%%", frac*100) }

// PadRight pads a string on the right to the given width. The helper is used
// to align the label column of the text report.
func PadRight(s string, width int) string {
	for len(s) < width {
		s += " "
	}
	return s
}

// Rule renders a full-width separator line of the given character, matching
// the length of the section title it follows.
func Rule(ch byte, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}

// Labeled formats a "label : value" pair with a fixed label width. Keeping the
// alignment in one function makes the report rows consistent.
func Labeled(label, value string, width int) string {
	return PadRight(label, width) + ": " + value
}

// SignedDelta formats a signed difference with a leading plus or minus sign,
// which is used for deviations from nominal values.
func SignedDelta(v float64, format string) string {
	sign := ""
	if v > 0 {
		sign = "+"
	}
	return sign + fmt.Sprintf(format, v)
}
