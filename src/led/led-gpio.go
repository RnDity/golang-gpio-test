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
	"io/ioutil"
	"os"
	"strconv"
)

type gpioLed struct {
	pin        int
	is_on      bool
	value_file *os.File
}

func export_gpio(pin int) {
	gpio_dir := fmt.Sprintf("/sys/class/gpio/gpio%d", pin)
	if _, err := os.Stat(gpio_dir); os.IsNotExist(err) {
		checked_write("/sys/class/gpio/export", strconv.Itoa(pin))
		fmt.Printf("GPIO %d exported\n", pin)
	} else {
		fmt.Printf("GPIO %d already exported\n", pin)
	}
}

func unexport_gpio(pin int) {
	fmt.Printf("unexporting GPIO %d\n", pin)
	checked_write("/sys/class/gpio/unexport", strconv.Itoa(pin))
}

func checked_write(path string, value string) {
	if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
		panic(fmt.Sprintf("write of value %s to %s failed: %s",
			value, path, err.Error()))
	}
}

func NewGPIO(pin int) Led {
	l := &gpioLed{pin: pin}
	return l
}

func (l *gpioLed) get_path(prop string) string {
	return fmt.Sprintf("/sys/class/gpio/gpio%d/%s", l.pin, prop)
}

func (l *gpioLed) value(out string) {
	_, err := l.value_file.WriteString(out)
	if err != nil {
		panic(fmt.Sprintf("write to GPIO %s failed: %s",
			l.get_path("value"), err.Error()))
	}
}

func (l *gpioLed) On() {
	l.value("1")
	l.is_on = true
}

func (l *gpioLed) Off() {
	l.value("0")
	l.is_on = false
}

func (l *gpioLed) Toggle() {
	if l.is_on == true {
		l.Off()
	} else {
		l.On()
	}
}

func (l *gpioLed) Init() bool {
	export_gpio(l.pin)
	checked_write(l.get_path("direction"), "out")

	f, err := os.OpenFile(l.get_path("value"), os.O_WRONLY, 0)
	if err != nil {
		panic(fmt.Sprintf("failed to open GPIO file %s: %s",
			l.get_path("value"), err.Error()))
	}
	l.value_file = f

	l.Off()
	return true
}

func (l *gpioLed) Close() {
	if l.value_file != nil {
		l.value_file.Close()
	}
	unexport_gpio(l.pin)
}
