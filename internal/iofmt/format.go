package iofmt

import "fmt"

func F4(v float64) string { return fmt.Sprintf("%.4f", v) }

func F2(v float64) string { return fmt.Sprintf("%.2f", v) }

func F6(v float64) string { return fmt.Sprintf("%.6f", v) }

func FG(v float64) string { return fmt.Sprintf("%.6g", v) }

func Pct1(frac float64) string { return fmt.Sprintf("%.1f%%", frac*100) }

func Pct2(frac float64) string { return fmt.Sprintf("%.2f%%", frac*100) }

func PadRight(s string, width int) string {
	for len(s) < width {
		s += " "
	}
	return s
}

func Rule(ch byte, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}

func Labeled(label, value string, width int) string {
	return PadRight(label, width) + ": " + value
}

func SignedDelta(v float64, format string) string {
	sign := ""
	if v > 0 {
		sign = "+"
	}
	return sign + fmt.Sprintf(format, v)
}
