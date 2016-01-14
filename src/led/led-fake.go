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
