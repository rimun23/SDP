package main

import (
	"errors"
	"fmt"
	"io"
)

type NotificationChannel interface {
	Send(title, message string) error
}

var (
	ErrChannelUnavailable = errors.New("channel unavailable")
	ErrMessageRejected    = errors.New("message rejected")
)

type ChannelError struct {
	Channel string
	Kind    error
	Detail  string
}

func (e *ChannelError) Error() string {
	return fmt.Sprintf("%s: %v (%s)", e.Channel, e.Kind, e.Detail)
}
func (e *ChannelError) Unwrap() error { return e.Kind }

type ConsoleChannel struct {
	Out io.Writer
}

func (c ConsoleChannel) Send(title, message string) error {
	if _, err := fmt.Fprintf(c.Out, "[console] %s — %s\n", title, message); err != nil {
		return &ChannelError{"console", ErrChannelUnavailable, err.Error()}
	}
	return nil
}

type EmailChannel struct {
	To  string
	Out io.Writer
}

func (e EmailChannel) Send(title, message string) error {
	if e.To == "" {
		return &ChannelError{"email", ErrChannelUnavailable, "no recipient configured"}
	}
	if _, err := fmt.Fprintf(e.Out, "[email to %s] %s\n%s\n", e.To, title, message); err != nil {
		return &ChannelError{"email", ErrChannelUnavailable, err.Error()}
	}
	return nil
}
