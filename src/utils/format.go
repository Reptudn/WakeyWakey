package utils

import "regexp"

var macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)

func IsValidMacAddress(mac string) bool {
	return macRegex.MatchString(mac)
}