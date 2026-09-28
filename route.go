package main

type Reading struct {
	Plant    string
	Moisture int
}
type Channels struct {
	Console NotificationChannel
	Email   NotificationChannel
	Buzzer  NotificationChannel
}

func Route(r Reading, ch Channels) Alert {
	switch {
	case r.Moisture < 15:
		return NewCriticalAlert(ch.Buzzer)
	case r.Moisture < 35:
		return NewThirstyAlert(ch.Email)
	case r.Moisture < 50:
		return NewThirstyAlert(ch.Console)
	default:
		return nil
	}
}
