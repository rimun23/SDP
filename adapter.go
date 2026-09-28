package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type buzzer interface {
	Buzz(code int, text string) int
}
type BuzzerAdapter struct {
	Device buzzer
}

const (
	buzzerBeeps   = 1
	buzzerMaxText = 32
)

func fitBuzzerText(title, message string) string {
	text := strings.TrimSpace(title + " " + message)

	if len(text) <= buzzerMaxText {
		return text
	}

	const ellipsis = "…"
	limit := buzzerMaxText - len(ellipsis)

	var result strings.Builder

	for _, r := range text {
		if result.Len()+utf8.RuneLen(r) > limit {
			break
		}
		result.WriteRune(r)
	}

	result.WriteString(ellipsis)
	return result.String()
}

func (a BuzzerAdapter) Send(title, message string) error {
	text := fitBuzzerText(title, message)
	code := a.Device.Buzz(buzzerBeeps, text)
	switch code {
	case BuzzOK:
		return nil
	case BuzzErrNotConnected:
		return &ChannelError{"buzzer", ErrChannelUnavailable, "device not connected"}

	case BuzzErrTextTooLong, BuzzErrBadCode:
		return &ChannelError{"buzzer", ErrMessageRejected, "device refused the message"}
	default:
		return &ChannelError{"buzzer", ErrChannelUnavailable, fmt.Sprintf("unknown device status %d", code)}
	}
}
