package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrWriteRune     = errors.New("error write rune")
)

func Unpack(s string) (string, error) {
	var (
		res      strings.Builder
		lastRune rune
	)

	for _, r := range s {
		if count, err := strconv.Atoi(string(r)); err == nil {
			if lastRune == 0 || unicode.IsDigit(lastRune) {
				return "", ErrInvalidString
			}

			for ; count > 1; count-- {
				if _, err := res.WriteRune(lastRune); err != nil {
					return "", ErrWriteRune
				}
			}
		} else if _, err := res.WriteRune(r); err != nil {
			return "", ErrWriteRune
		}

		lastRune = r
	}

	return res.String(), nil
}
