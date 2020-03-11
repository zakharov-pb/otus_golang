package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

// ErrInvalidString error returned from function Unpack
var ErrInvalidString = errors.New("invalid string")

// Unpack function to unzip an archived string
func Unpack(archive string) (string, error) {
	var result strings.Builder
	var prevSymb rune = -1
	shielding := false
	for _, r := range archive {
		if unicode.IsDigit(r) && !shielding {
			if prevSymb < 0 {
				return "", ErrInvalidString
			}
			count := int(r - '0')
			for i := 0; i < count; i++ {
				result.WriteRune(prevSymb)
			}
			prevSymb = -1
		} else {
			if shielding && !unicode.IsDigit(r) && r != '\\' {
				return "", ErrInvalidString
			}
			if !shielding && r == '\\' {
				shielding = true
				continue
			}

			if prevSymb > 0 {
				result.WriteRune(prevSymb)
			}
			prevSymb = r
			shielding = false
		}
	}
	if prevSymb > 0 {
		result.WriteRune(prevSymb)
	}
	return result.String(), nil
}
