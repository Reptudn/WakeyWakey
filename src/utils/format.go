package utils

func IsValidMacAddress(mac string) bool {
	if len(mac) != 17 {
		return false
	}

	for i, c := range mac {
		if i%3 == 2 {
			if c != ':' && c != '-' {
				return false
			}
		} else {
			if (c < '0' || c > '9') && (c < 'A' || c > 'F') && (c < 'a' || c > 'f') {
				return false
			}
		}
	}

	return true
}