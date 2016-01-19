// LED definitions for Raspberry PI
package led

func GetRaspberryPiLEDs() []Led {
	return NewForGPIOPins([]int{4, 17, 27, 22})
}
