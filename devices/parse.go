package devices

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dermaddis/z2m_cli/sliceutil"
	"github.com/hashicorp/go-multierror"
)

type Device string

var deviceNames = []string{
	"strip01",
	"light01",
}

var deviceNameAliases = map[string]string{
	"s01": "strip01",
	"s1":  "strip01",
	"l01": "light01",
	"l1":  "light01",
}

var UnknownDevice = errors.New("unknown device")

func Parse(s string) ([]Device, error) {
	if len(s) == 0 {
		return []Device{}, errors.New("s cannot be empty")
	}

	s = strings.ToLower(s)
	split := strings.Split(s, ",")

	var allErrors error
	devices := make([]Device, 0, len(deviceNames))

	for _, deviceString := range split {
		if deviceString == "all" {
			// We put all deviceNames into devices
			for i, deviceName := range deviceNames {
				if len(devices)-1 < i {
					// This index is not yet used => append
					devices = append(devices, Device(deviceName))
				} else {
					// This index was already used => overwrite
					devices[i] = Device(deviceName)
				}
			}
			break
		}

        // Search for alias and use it if exists
		deviceName, ok := deviceNameAliases[deviceString]
		if ok {
			deviceString = deviceName
		}

		contains := sliceutil.Contains(deviceNames, deviceString)
		if contains {
			devices = append(devices, Device(deviceString))
		} else {
			allErrors = multierror.Append(allErrors, fmt.Errorf("%q: %w", deviceString, UnknownDevice))
		}
	}

	return devices, allErrors
}
