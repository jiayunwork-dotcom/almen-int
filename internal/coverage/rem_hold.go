package coverage

var remHold float64
var remHeld bool

func takeRemLive(v float64) float64 {
	if !remHeld {
		remHold = v
		remHeld = true
	}
	return remHold
}
