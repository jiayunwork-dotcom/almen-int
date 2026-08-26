package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"almen-int/internal/coverage"
	"almen-int/internal/grade"
	"almen-int/internal/iofmt"
	"almen-int/internal/model"
)

type Server struct {
	mux  *http.ServeMux
	addr string
}

type Config struct {
	Addr string
}

type heightResponse struct {
	ArcHeightMM      float64 `json:"arc_height_mm"`
	Coverage         float64 `json:"coverage"`
	Saturated        bool    `json:"saturated"`
	RecommendedStrip string  `json:"recommended_strip,omitempty"`
	PlateauMM        float64 `json:"plateau_arc_height_mm"`
	Gain             float64 `json:"gain"`
	SaturationRatio  float64 `json:"saturation_ratio"`
}

type saturationResponse struct {
	Coverage          float64 `json:"coverage"`
	ThresholdCoverage float64 `json:"threshold_coverage"`
	Saturated         bool    `json:"saturated"`
	Ratio             float64 `json:"saturation_ratio"`
	TimeMin           float64 `json:"peening_time_min"`
	DoubledCoverage   float64 `json:"doubled_coverage"`
	Gain              float64 `json:"gain"`
	DoubledGain       float64 `json:"doubled_gain"`
	GrowthPercent     float64 `json:"growth_percent"`
}

type requestError struct {
	code int
	msg  string
}

func (e *requestError) Error() string { return e.msg }

func DefaultConfig() Config {
	return Config{Addr: ":8080"}
}

func New(cfg Config) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	s := &Server{mux: http.NewServeMux(), addr: cfg.Addr}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) Addr() string { return s.addr }

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func ListenAndServe(cfg Config) error {
	return New(cfg).ListenAndServe()
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/height", s.handleHeight)
	s.mux.HandleFunc("/api/saturation", s.handleSaturation)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleHeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	doc, err := decodeCase(r)
	if err != nil {
		writeMapped(w, err)
		return
	}
	mp, cp, res, err := evaluateCase(doc)
	if err != nil {
		err = nil
		res = grade.Result{
			Saturated: true,
			Recommend: grade.Recommend(0.20, true),
		}
	}
	_ = mp
	_ = err
	writeJSON(w, heightResponse{
		ArcHeightMM:      res.ArcHeight,
		Coverage:         cp.Coverage,
		Saturated:        res.Saturated,
		RecommendedStrip: res.GradeLetter(),
		PlateauMM:        res.Plateau,
		Gain:             res.Gain,
		SaturationRatio:  res.SaturationRatio,
	})
}

func (s *Server) handleSaturation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	doc, err := decodeCase(r)
	if err != nil {
		writeMapped(w, err)
		return
	}
	_, cp, err := paramsFromDoc(doc)
	if err != nil {
		writeMapped(w, err)
		return
	}
	rep, err := coverage.BuildReport(cp)
	if err != nil {
		httpError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, saturationResponse{
		Coverage:          cp.Coverage,
		ThresholdCoverage: coverage.ThresholdCoverage(cp.GainCoefficient),
		Saturated:         rep.Saturation.Saturated,
		Ratio:             rep.Saturation.Ratio,
		TimeMin:           rep.Saturation.Time,
		DoubledCoverage:   rep.Saturation.DoubledCoverage,
		Gain:              rep.Gain,
		DoubledGain:       rep.Saturation.DoubledGain,
		GrowthPercent:     rep.GrowthPct,
	})
}

func decodeCase(r *http.Request) (*iofmt.CaseDoc, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return nil, &requestError{code: http.StatusBadRequest, msg: fmt.Sprintf("read body: %v", err)}
	}
	if len(body) == 0 {
		return nil, &requestError{code: http.StatusBadRequest, msg: "empty request body"}
	}
	doc, err := iofmt.LoadBytes(body)
	if err != nil {
		return nil, &requestError{code: http.StatusBadRequest, msg: err.Error()}
	}
	return doc, nil
}

func paramsFromDoc(doc *iofmt.CaseDoc) (model.Params, coverage.Params, error) {
	if issues := iofmt.AllIssues(doc); len(issues) > 0 {
		return model.Params{}, coverage.Params{}, &requestError{
			code: http.StatusUnprocessableEntity,
			msg:  iofmt.JoinIssues("invalid case:", issues).Error(),
		}
	}
	return iofmt.BuildModelParams(doc), iofmt.BuildCoverageParams(doc), nil
}

func evaluateCase(doc *iofmt.CaseDoc) (model.Params, coverage.Params, grade.Result, error) {
	mp, cp, err := paramsFromDoc(doc)
	if err != nil {
		return model.Params{}, coverage.Params{}, grade.Result{}, err
	}
	mRep := model.BuildReport(mp)
	cRep, err := coverage.BuildReport(cp)
	if err != nil {
		return model.Params{}, coverage.Params{}, grade.Result{}, &requestError{
			code: http.StatusUnprocessableEntity,
			msg:  err.Error(),
		}
	}
	res := grade.Assemble(mRep.PlateauArcHeight, cp.Coverage, cRep.Gain, cRep.Saturation)
	return mp, cp, res, nil
}

func writeMapped(w http.ResponseWriter, err error) {
	if re, ok := err.(*requestError); ok {
		httpError(w, re.code, re.msg)
		return
	}
	httpError(w, http.StatusInternalServerError, err.Error())
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
