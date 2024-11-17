package assert

import (
	"log"
)

func Nil(value any, msgs ...any) {
	if value != nil {
		log.Panicln(append([]any{"Assertion Failed - Value must be 'nil' but received:", value}, resolveMsgs(msgs)...)...)
	}
}

func True(condition bool, msgs ...any) {
	if !condition {
		log.Panicln(append([]any{"Assertion Failed!"}, resolveMsgs(msgs)...)...)
	}
}

func False(condition bool, msgs ...any) {
	if condition {
		log.Panicln(append([]any{"Assertion Failed!"}, resolveMsgs(msgs)...)...)
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
