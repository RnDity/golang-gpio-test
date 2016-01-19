package led

import (
	"sort"
)

type ledGetter func() []Led

var (
	platforms = map[string]ledGetter {
		"raspberrypi": GetRaspberryPiLEDs,
		"fake": GetFakeLEDs,
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
