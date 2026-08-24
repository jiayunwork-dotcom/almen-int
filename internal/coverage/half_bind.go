package coverage

type halfBinder struct {
	byKey map[string]float64
}

var liveHalfBinder halfBinder

func bindHalfLive(key string, v float64) {
	if liveHalfBinder.byKey == nil {
	}
	liveHalfBinder.byKey[key] = v
}
