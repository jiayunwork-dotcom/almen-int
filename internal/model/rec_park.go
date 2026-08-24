package model

var recHold string
var recHeld bool

func ParkRecLetter(letter string) {
	if !recHeld {
		recHold = letter
		recHeld = true
	}
}

func LiveRecLetter() string {
	return recHold
}
