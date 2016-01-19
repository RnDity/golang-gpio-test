package led

import (
	"sort"
)

type ledGetter func() []Led

var (
	platforms = map[string]ledGetter{
		"raspberrypi": GetRaspberryPiLEDs,
		"beaglebone":  GetBeagleBoneBlackLEDs,
		"fake":        GetFakeLEDs,
		"olimexa20":   GetOlimexA20LEDs,
	}
)

func ListPlatforms() []string {
	var keys []string
	for k := range platforms {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func LEDsForPlatform(plat string) []Led {
	getter, exist := platforms[plat]
	if exist {
		return getter()
	}
	return []Led{}
}
