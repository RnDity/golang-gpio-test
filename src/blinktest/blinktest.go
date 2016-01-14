package main

import (
	"chain"
	"fmt"
	"led"
	"time"
)

const (
	REVOLUTIONS = 100
)

func main() {
	leds := []led.Led{
		led.NewGPIO(4),
		led.NewGPIO(17),
		led.NewGPIO(27),
		led.NewGPIO(22),
	}
	// leds := []led.Led{
	// 	led.NewFake(4),
	// 	led.NewFake(17),
	// 	led.NewFake(24),
	// 	led.NewFake(22),
	// }

	for _, l := range leds {
		l.Init()
	}

	chains := []*chain.Chain{
		chain.New(REVOLUTIONS, 1000, leds[0]),
		chain.New(REVOLUTIONS, 10000, leds[1]),
		chain.New(REVOLUTIONS, 25000, leds[2]),
		chain.New(REVOLUTIONS, 100000, leds[3]),
	}

	fmt.Println("spawning goroutines")
	// spawn all goroutines
	for _, ch := range chains {
		ch.Spawn()
	}

	time.Sleep(time.Second)
	// start all
	fmt.Println("trigger start")
	for _, ch := range chains {
		ch.Start()
	}

	for _, ch := range chains {
		ch.Wait()
	}
	fmt.Println("finished")

	for _, l := range leds {
		l.Close()
	}
}
