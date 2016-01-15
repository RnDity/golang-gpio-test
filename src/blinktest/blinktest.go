package main

import (
	"chain"
	"flag"
	"fmt"
	"led"
	"os"
	"time"
)

const (
	REVOLUTIONS = 100
	REV_PROC = 100000 * REVOLUTIONS

	usage = "" +
		`Usage of 'blinktest' command
Flags:
	-usefake=[0|1]:   Use fake LEDs
`
)

var (
	use_fake  = flag.Bool("usefake", false, "Use fake LEDs")

	// Raspberry PI LEDs connected via GPIO
	gpio_leds = []led.Led{
		led.NewGPIO(4),
		led.NewGPIO(17),
		led.NewGPIO(27),
		led.NewGPIO(22),
	}

	fake_leds = []led.Led{
		led.NewFake(4),
		led.NewFake(17),
		led.NewFake(24),
		led.NewFake(22),
	}
)

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
		os.Exit(1)
	}

	flag.Parse()

	var leds []led.Led

	if *use_fake {
		fmt.Println("using fake LEDs")
		leds = fake_leds
	} else {
		leds = gpio_leds
	}

	for _, l := range leds {
		l.Init()
	}

	chains := []*chain.Chain{
		chain.New(REV_PROC / 1000, 1000, leds[0]),
		chain.New(REV_PROC / 10000, 10000, leds[1]),
		chain.New(REV_PROC / 25000, 25000, leds[2]),
		chain.New(REV_PROC / 100000, 100000, leds[3]),
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
