package assert

import (
	"fmt"
	"log"
)

func Nil(value any, msgs ...any) {
	if value != nil {
		log.Panicln(append([]any{"Assertion Failed - Value must be 'nil' but received:", value}, resolveMsgs(msgs)...)...)
	}
}

func Eq(a any, b any, msgs ...any) {
	if a != b {
		log.Panicln(append([]any{fmt.Sprintf("Assertion Failed - 'a' must be equal to 'b'. Received\na: %s\n\nb: %s", a, b)}, resolveMsgs(msgs)...)...)
	}
}

func Neq(a any, b any, msgs ...any) {
	if a == b {
		log.Panicln(append([]any{fmt.Sprintf("Assertion Failed - 'a' must not be equal to 'b'. Received\na & b: %s", a)}, resolveMsgs(msgs)...)...)
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
