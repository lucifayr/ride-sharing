package assert

import (
	"log"
)

func True(condition bool, msgs ...any) {
	if !condition {
		log.Fatalln(append([]any{"Assertion Failed!"}, resolveMsgs(msgs)...)...)
	}
}

func False(condition bool, msgs ...any) {
	if condition {
		log.Fatalln(append([]any{"Assertion Failed!"}, resolveMsgs(msgs)...)...)
	}
}

func resolveMsgs(msgs []any) []any {
	msgsResolved := make([]any, len(msgs))
	for idx, msg := range msgs {
		f, ok := msg.(func() any)
		if ok {
			msgsResolved[idx] = f()
		} else {
			msgsResolved[idx] = msg
		}
	}

	return msgsResolved
}
