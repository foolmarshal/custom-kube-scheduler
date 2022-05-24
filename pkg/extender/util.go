package extender

import (
	"github.com/docker/go-units"
	"strconv"
	"strings"
)

// ConvertToIBytes humanizes all the passed units to IBytes format
func ConvertToIBytes(value string) string {
	if value == "" {
		return value
	}
	var bytes int64
	var err error
	if strings.Contains(value, "i") {
		bytes, err = units.RAMInBytes(value)
	} else {
		bytes, err = units.FromHumanSize(value)
	}
	// should error be also returned?
	if err != nil {
		return value
	}
	return units.CustomSize("%.1f%s", float64(bytes), 1024.0, []string{"B"})
}

func extractIntegralPartFromVolumeInBytes(value string) (int64, error) {
	bytes, err := strconv.Atoi(strings.Split(strings.Split(value, "B")[0], ".")[0])
	return int64(bytes), err
}
