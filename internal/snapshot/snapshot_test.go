package snapshot

import (
	"os"
	"path/filepath"
	"testing"

	"almen-int/internal/coverage"
	"almen-int/internal/model"
)

func sampleModel() model.Params {
	return model.Params{
		Velocity:          48.0,
		ReferenceVelocity: 50.0,
		ShotDiameter:      0.6,
		ShotDensity:       7800.0,
		Thickness:         1.29,
		Width:             18.5,
		Length:            76.0,
		Modulus:           205.0,
		ResidualStress:    850.0,
		LayerDepth:        0.05,
	}
}

func sampleCoverage() coverage.Params {
	return coverage.Params{
		Coverage:        0.98,
		RateConstant:    0.085,
		GainCoefficient: 2.6,
	}
}

func TestSnapshotRoundTripAgrees(t *testing.T) {
	rec, err := Capture(sampleModel(), sampleCoverage())
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "almen.snap.json")
	if err := WriteFile(path, rec); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Matches(rec) {
		t.Fatalf("round-trip mismatch: %+v vs %+v", got, rec)
	}
	if err := got.ReplayAgrees(); err != nil {
		t.Fatal(err)
	}
	if got.GradeLetter != "A" {
		t.Fatalf("saturated A-strip case must store letter A, got %q", got.GradeLetter)
	}
}

func TestSnapshotTruncationKeepsPriorFile(t *testing.T) {
	rec, err := Capture(sampleModel(), sampleCoverage())
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	good := filepath.Join(dir, "good.json")
	if err := WriteFile(good, rec); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "trunc.json")
	raw, err := os.ReadFile(good)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) < 40 {
		t.Fatal("snapshot too small to truncate")
	}
	if err := os.WriteFile(bad, raw[:len(raw)/2], 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(bad); err == nil {
		t.Fatal("truncated JSON must be rejected")
	}
	kept, err := ReadFile(good)
	if err != nil {
		t.Fatal(err)
	}
	if err := kept.ReplayAgrees(); err != nil {
		t.Fatal(err)
	}
	if !kept.Matches(rec) {
		t.Fatal("prior snapshot must still match the live kernel")
	}
}

func TestEmptyFileRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path); err == nil {
		t.Fatal("empty file must be rejected")
	}
}

func TestCaptureRejectsIllegalModel(t *testing.T) {
	mp := sampleModel()
	mp.Velocity = 0
	if _, err := Capture(mp, sampleCoverage()); err == nil {
		t.Fatal("illegal velocity must not snapshot")
	}
}

func TestUnsaturatedSnapshotHasNoGrade(t *testing.T) {
	cp := sampleCoverage()
	cp.Coverage = 0.5
	rec, err := Capture(sampleModel(), cp)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Saturated {
		t.Fatal("coverage 0.5 must not be saturated")
	}
	if rec.GradeLetter != "" {
		t.Fatalf("unsaturated snapshot must not store a grade, got %q", rec.GradeLetter)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "unsat.json")
	if err := WriteFile(path, rec); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := got.ReplayAgrees(); err != nil {
		t.Fatal(err)
	}
}

func TestTamperedArcHeightFailsReplay(t *testing.T) {
	rec, err := Capture(sampleModel(), sampleCoverage())
	if err != nil {
		t.Fatal(err)
	}
	rec.ArcHeightMM = rec.ArcHeightMM * 1.5
	if err := rec.ReplayAgrees(); err == nil {
		t.Fatal("tampered arc height must fail replay")
	}
}
