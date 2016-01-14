package main

import (
	"fmt"
	"led"
	"time"
)

func main() {
	leds := []led.Led{
		led.NewGPIO(4),
		led.NewGPIO(17),
		led.NewGPIO(27),
		led.NewGPIO(22),
	}

	for _, l := range leds {
		l.Init()
	}

	for {
		fmt.Println("toggle")
		for _, l := range leds {
			l.Toggle()
		}

		time.Sleep(time.Second)
	}
}
