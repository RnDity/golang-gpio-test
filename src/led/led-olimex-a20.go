package led

import (
	"fmt"
	"os"
	"strconv"
)

type olimexa20Led struct {
	gpioLed
	signalName string
}

func NewOlimexGPIO(pin int, signal string) Led {
	ol := &olimexa20Led{
		gpioLed: gpioLed{
			pin: pin,
		},
		signalName: signal,
	}

	return ol
}

func (l *olimexa20Led) get_path(prop string) string {
	return fmt.Sprintf("/sys/class/gpio/gpio%d_%s/%s",
		l.pin, l.signalName, prop)
}

func export_olimex_gpio(pin int, signal string) {
	gpio_dir := fmt.Sprintf("/sys/class/gpio/gpio%d_%s", pin, signal)
	if _, err := os.Stat(gpio_dir); os.IsNotExist(err) {
		checked_write("/sys/class/gpio/export", strconv.Itoa(pin))
		fmt.Printf("GPIO %d exported\n", pin)
	} else {
		fmt.Printf("GPIO %d already exported\n", pin)
	}
}

func (l *olimexa20Led) Init() bool {
	export_olimex_gpio(l.pin, l.signalName)
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

func GetOlimexA20LEDs() []Led {
	return []Led{
		NewOlimexGPIO(35, "pb3"),
		NewOlimexGPIO(37, "pb5"),
		NewOlimexGPIO(39, "pb7"),
		NewOlimexGPIO(40, "pb10"),
	}
}
