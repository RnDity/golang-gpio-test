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

package led

import (
	"fmt"
)

const (
	OFF = 0
	ON  = 1
)

type ledstate int

type fakeLed struct {
	id    int
	state ledstate
}

func NewFake(id int) Led {
	return &fakeLed{id, OFF}
}

func (l fakeLed) report_state() {
	state := "ON"
	if l.state == OFF {
		state = "OFF"
	}
	fmt.Printf("[%d] state: %s\n", l.id, state)
}

func (l *fakeLed) set_state(state ledstate) {
	l.state = state
	l.report_state()
}

func (l *fakeLed) On() {
	l.set_state(ON)
}

func (l *fakeLed) Off() {
	l.set_state(OFF)
}

func (l *fakeLed) Toggle() {
	if l.state == ON {
		l.Off()
	} else {
		l.On()
	}
}

func (l *fakeLed) Init() bool {
	l.Off()
	return true
}

func (l *fakeLed) Close() {

}

func GetFakeLEDs() []Led {
	return []Led{
		NewFake(4),
		NewFake(17),
		NewFake(24),
		NewFake(22),
	}
}
