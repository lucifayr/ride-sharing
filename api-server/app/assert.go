package assert

import "log"

func True(condition bool, msgs ...any) {
	if !condition {
		log.Fatalln("Assertion Failed!", msgs)
	}
}

func False(condition bool, msgs ...any) {
	if condition {
		log.Fatalln("Assertion Failed!", msgs)
	}
}
