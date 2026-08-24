package iofmt

import "testing"

func TestFormatHelpers(t *testing.T) {
	if got := F4(0.44068); got != "0.4407" {
		t.Errorf("F4 must round to four decimals, got %q", got)
	}
	if got := Pct1(0.98); got != "98.0%" {
		t.Errorf("Pct1 must render 98.0%%, got %q", got)
	}
	if got := Labeled("arc height", "0.44 mm", 18); got != "arc height        : 0.44 mm" {
		t.Errorf("Labeled must align the label column, got %q", got)
	}
	if got := Rule('-', 4); got != "----" {
		t.Errorf("Rule must repeat the character, got %q", got)
	}
	if got := SignedDelta(0.02, "%.2f"); got != "+0.02" {
		t.Errorf("SignedDelta must prefix a plus sign, got %q", got)
	}
	if got := SignedDelta(-0.02, "%.2f"); got != "-0.02" {
		t.Errorf("SignedDelta must keep the minus sign, got %q", got)
	}
}

func TestDescribeCase(t *testing.T) {
	doc, err := LoadBytes([]byte(validCase))
	if err != nil {
		t.Fatalf("LoadBytes failed: %v", err)
	}
	summary := DescribeCase(doc)
	for _, want := range []string{"A2 tool steel", "v=48.0 m/s", "t=1.29 mm", "E=205 GPa", "coverage=0.98"} {
		if !contains(summary, want) {
			t.Errorf("case summary must mention %q, got %q", want, summary)
		}
	}
}

func TestDescribeParams(t *testing.T) {
	doc, _ := LoadBytes([]byte(validCase))
	mp := BuildModelParams(doc)
	s := DescribeParams(mp)
	if !contains(s, "1.29 mm") || !contains(s, "18.5 mm") || !contains(s, "76.0 mm") {
		t.Errorf("params description must carry the strip size, got %q", s)
	}
}

func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
