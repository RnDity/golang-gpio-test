// Copyright (c) 2015 Open-RnD Sp. z o.o.

// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"chain"
	"flag"
	"fmt"
	"led"
	"os"
	"time"
)

type procSetConf struct {
	revs  int // complete chain passes
	count int // processes count
}

const (
	REVOLUTIONS = 100
	REV_PROC    = 100000 * REVOLUTIONS

	usage = "" +
		`Usage of 'blinktest' command
Flags:
	-usefake:    Use fake LEDs
	-continuous: Toggle GPIO lines for approximately the same time
`
)

var (
	use_fake   = flag.Bool("usefake", false, "Use fake LEDs")
	continuous = flag.Bool("continuous", false, "Use fake LEDs")

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
	procsetOneByOne = []procSetConf{
		{REVOLUTIONS, 1000},
		{REVOLUTIONS, 10000},
		{REVOLUTIONS, 25000},
		{REVOLUTIONS, 100000},
	}

	procsetContinuous = []procSetConf{
		{REV_PROC / 1000, 1000},
		{REV_PROC / 10000, 10000},
		{REV_PROC / 25000, 25000},
		{REV_PROC / 100000, 100000},
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

	procset := procsetOneByOne
	if *continuous {
		fmt.Println("using continous processes set")
		procset = procsetContinuous
	}

	if len(procset) > len(leds) {
		fmt.Fprintf(os.Stderr, "Processes set larger than avaialble GPIOs")
		os.Exit(1)
	}

	// setup process chains
	chains := make([]*chain.Chain, len(procset))
	for i := range chains {
		fmt.Printf("adding chain %d: %d passes %d processes\n",
			i, procset[i].revs, procset[i].count)
		chains[i] = chain.New(procset[i].revs, procset[i].count, leds[i])
	}

	for _, l := range leds {
		l.Init()
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

	// wait for chains to complete
	for _, ch := range chains {
		ch.Wait()
	}
	fmt.Println("finished")

	// cleanup LEDs
	for _, l := range leds {
		l.Close()
	}
}
