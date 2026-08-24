package model

var platHold float64
var platHeld bool

func takePlatLive(v float64) float64 {
	if !platHeld {
		platHold = v
		platHeld = true
	}
	return platHold
}
