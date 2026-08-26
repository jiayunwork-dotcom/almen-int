package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const validCase = `{
  "shot": {"velocity": 48.0, "diameter": 0.60, "density": 7800.0},
  "strip": {"thickness": 1.29, "width": 18.5, "length": 76.0, "modulus": 205.0},
  "residual": {"stress": 850.0, "layer_depth": 0.05},
  "process": {"coverage": 0.98, "rate_constant": 0.085, "gain_coefficient": 2.6},
  "reference": {"velocity": 50.0}
}`

func TestHealth(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("health code %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "ok") {
		t.Fatalf("health body %q", rec.Body.String())
	}
}

func TestHeightEndpoint(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/height", strings.NewReader(validCase))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("height code %d body %s", rec.Code, rec.Body.String())
	}
	var resp heightResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ArcHeightMM <= 0 {
		t.Fatalf("arc_height_mm must be positive, got %v", resp.ArcHeightMM)
	}
	if !resp.Saturated {
		t.Fatalf("reference coverage must be saturated")
	}
	if resp.RecommendedStrip != "A" {
		t.Fatalf("expected strip A, got %q", resp.RecommendedStrip)
	}
}

func TestHeightIllegalVelocity(t *testing.T) {
	srv := New(DefaultConfig())
	raw := []byte(`{
	  "shot": {"velocity": 0, "diameter": 0.60, "density": 7800.0},
	  "strip": {"thickness": 1.29, "width": 18.5, "length": 76.0, "modulus": 205.0},
	  "residual": {"stress": 850.0, "layer_depth": 0.05},
	  "process": {"coverage": 0.98, "rate_constant": 0.085, "gain_coefficient": 2.6},
	  "reference": {"velocity": 50.0}
	}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/height", bytes.NewReader(raw))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "error") {
		t.Fatalf("body should contain error: %q", rec.Body.String())
	}
}

func TestHeightInvalidJSON(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/height", strings.NewReader(`{bad`))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHeightMethodNotAllowed(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/height", nil)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestSaturationEndpoint(t *testing.T) {
	srv := New(DefaultConfig())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/saturation", strings.NewReader(validCase))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("saturation code %d body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "saturation_ratio") {
		t.Fatalf("expected saturation_ratio in body, got %s", rec.Body.String())
	}
}

func TestSaturationIllegalCoverage(t *testing.T) {
	srv := New(DefaultConfig())
	raw := []byte(`{
	  "shot": {"velocity": 48.0, "diameter": 0.60, "density": 7800.0},
	  "strip": {"thickness": 1.29, "width": 18.5, "length": 76.0, "modulus": 205.0},
	  "residual": {"stress": 850.0, "layer_depth": 0.05},
	  "process": {"coverage": 0, "rate_constant": 0.085, "gain_coefficient": 2.6},
	  "reference": {"velocity": 50.0}
	}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/saturation", bytes.NewReader(raw))
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d body %s", rec.Code, rec.Body.String())
	}
}
