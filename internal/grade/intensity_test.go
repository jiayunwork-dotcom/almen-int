package grade

import "testing"

func TestIntensityNotation(t *testing.T) {
	if got := IntensityNotation(0.4407, "A"); got != "0.4407 mmA" {
		t.Errorf("intensity notation must be 0.4407 mmA, got %q", got)
	}
	if got := IntensityNotation(0.23, "N"); got != "0.2300 mmN" {
		t.Errorf("intensity notation must be 0.2300 mmN, got %q", got)
	}
}

func TestParseIntensityNotation(t *testing.T) {
	val, letter, err := ParseIntensityNotation("0.4407 mmA")
	if err != nil {
		t.Fatalf("ParseIntensityNotation must accept a valid notation: %v", err)
	}
	if !closeEnough(val, 0.4407, 1e-12) {
		t.Errorf("parsed value must be 0.4407, got %g", val)
	}
	if letter != "A" {
		t.Errorf("parsed letter must be A, got %q", letter)
	}

	if _, _, err := ParseIntensityNotation("0.44"); err == nil {
		t.Errorf("a notation without the mm unit must be rejected")
	}
	if _, _, err := ParseIntensityNotation("abc mmA"); err == nil {
		t.Errorf("a non-numeric arc height must be rejected")
	}
	if _, _, err := ParseIntensityNotation(""); err == nil {
		t.Errorf("an empty notation must be rejected")
	}
}

func TestIntensityRoundTrip(t *testing.T) {
	for _, letter := range []string{"N", "A", "C"} {
		if !RoundTripIntensity(0.1234, letter) {
			t.Errorf("round trip must reproduce the notation for %q", letter)
		}
	}
}

func TestIntensityTextGatedOnSaturation(t *testing.T) {
	saturated := Recommend(0.44, true)
	rSat := Result{ArcHeight: 0.44, Saturated: true, Recommend: saturated}
	if got := rSat.IntensityText(); got != "0.4400 mmA" {
		t.Errorf("saturated intensity must be 0.4400 mmA, got %q", got)
	}

	notSat := Recommend(0.44, false)
	rUns := Result{ArcHeight: 0.44, Saturated: false, Recommend: notSat}
	if got := rUns.IntensityText(); got != "-" {
		t.Errorf("unsaturated intensity must be a dash, got %q", got)
	}
}

func TestValidLetter(t *testing.T) {
	for _, letter := range []string{"N", "A", "C"} {
		if !ValidLetter(letter) {
			t.Errorf("%q must be a valid strip letter", letter)
		}
	}
	if ValidLetter("X") {
		t.Errorf("an unknown letter must be rejected")
	}
}
