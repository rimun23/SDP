package main

import (
	"strings"
)

const (
	BuzzOK              = 0
	BuzzErrNotConnected = -1
	BuzzErrTextTooLong  = -2
	BuzzErrBadCode      = -3
)

type BuzzerDevice struct {
	Connected bool
	Log       []string
}

func (b *BuzzerDevice) Buzz(code int, text string) int {
	if !b.Connected {
		return BuzzErrNotConnected
	}
	if code < 1 || code > 5 {
		return BuzzErrBadCode
	}
	if len(text) > 32 {
		return BuzzErrTextTooLong
	}
	b.Log = append(b.Log, strings.Repeat("BEEP ", code)+text)
	return BuzzOK
}
