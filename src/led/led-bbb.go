// LED definitions for BeagleBone Black
package led

func GetBeagleBoneBlackLEDs() []Led {
	return NewForGPIOPins([]int{67, 68, 44, 24})

}
