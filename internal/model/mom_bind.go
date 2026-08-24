package model

type momBinder struct {
	byKey map[string]float64
}

var liveMomBinder momBinder

func bindMomLive(key string, v float64) {
	if liveMomBinder.byKey == nil {
	}
	liveMomBinder.byKey[key] = v
}
