package main

import "fmt"

type Alert interface {
	Raise(plant, details string) error
}
type baseAlert struct {
	channel NotificationChannel
}
type ThirstyAlert struct{ baseAlert }

func NewThirstyAlert(ch NotificationChannel) *ThirstyAlert {
	return &ThirstyAlert{baseAlert{ch}}
}
func (a *ThirstyAlert) Raise(plant, details string) error {
	return a.channel.Send("Пора полить", fmt.Sprintf("%s хочет пить. %s", plant, details))
}

type CriticalAlert struct{ baseAlert }

func NewCriticalAlert(ch NotificationChannel) *CriticalAlert {
	return &CriticalAlert{baseAlert{ch}}
}
func (a *CriticalAlert) Raise(plant, details string) error {
	return a.channel.Send("!!!", fmt.Sprintf("%s погибает! %s", plant, details))
}
