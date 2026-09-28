package main

import (
	"fmt"
	"os"
)

func main() {
	buzzer := &BuzzerDevice{Connected: true}
	channels := Channels{
		Console: ConsoleChannel{Out: os.Stdout},
		Email:   EmailChannel{To: "example@example.com", Out: os.Stdout},
		Buzzer:  BuzzerAdapter{Device: buzzer},
	}

	readings := []Reading{
		{"Фикус", 45},
		{"Монстера", 30},
		{"Кактус", 14},
		{"Орхидея", 80},
	}

	for _, r := range readings {
		alert := Route(r, channels)

		if alert == nil {
			fmt.Printf("%s: всё в порядке\n", r.Plant)
			continue
		}

		previousBuzzerEntries := len(buzzer.Log)

		if err := alert.Raise(
			r.Plant,
			fmt.Sprintf("Влажность %d%%.", r.Moisture),
		); err != nil {
			fmt.Println("ошибка:", err)
			continue
		}

		// Выводим только новые записи устройства.
		if len(buzzer.Log) > previousBuzzerEntries {
			fmt.Println(
				"[buzzer]",
				buzzer.Log[len(buzzer.Log)-1],
			)
		}
	}
}
