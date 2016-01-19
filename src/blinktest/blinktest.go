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
		`Usage of 'blinktest' command:

	blinktest [flags] <platform-name>

Flags:
	-usefake:       Use fake LEDs
	-continuous:    Toggle GPIO lines for approximately the same time
	-useyield:      Enable voluntary CPU yielding by chain processes
	-listplatforms: List available platforms
`
)

var (
	use_fake       = flag.Bool("usefake", false, "Use fake LEDs")
	continuous     = flag.Bool("continuous", false, "Use fake LEDs")
	use_yield      = flag.Bool("useyield", false, "Enable voluntary CPU yielding")
	list_platforms = flag.Bool("listplatforms", false, "List avaialable platforms")

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

func listPlatforms() {
	fmt.Println("Available platforms:")
	for _, p := range led.ListPlatforms() {
		fmt.Printf("  %s\n", p)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
		os.Exit(1)
	}

	flag.Parse()

	if *list_platforms {
		listPlatforms()
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Println("Platform not specified, see -help for usage information")
		os.Exit(1)
	}

	platform := flag.Arg(0)

	if *use_yield {
		fmt.Println("enable voluntary yielding")
		chain.SetYield(true)
	}

	leds := led.LEDsForPlatform(platform)
	if len(leds) == 0 {
		fmt.Printf("no leds for platform '%s'\n", platform)
		os.Exit(1)
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
