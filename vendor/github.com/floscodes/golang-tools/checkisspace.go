package tools

import "unicode"

// Check if a string contains only space runes

func CheckIsSpace(s string) bool {
	space := true

	for _, x := range s {
		if !unicode.IsSpace(x) {
			space = false
			break
		}
	}

	return space
}
