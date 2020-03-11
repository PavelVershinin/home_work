package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const backslash = '\\'

var ErrInvalidString = errors.New("invalid string")
var ErrorOutOfRange = errors.New("out of range")

func Unpack(s string) (string, error) {
	var src = []rune(s)
	var dst strings.Builder

	if len(src) == 0 {
		return "", nil
	}

	for i := 0; i < len(src); i++ {
		symbol := src[i]
		count := 1

		if symbol == backslash {
			if escaped, err := RuneEscaped(src, i); err != nil {
				return "", err
			} else if !escaped {
				continue
			}
		}

		if unicode.IsDigit(symbol) {
			if escaped, err := RuneEscaped(src, i); err != nil {
				return "", err
			} else if !escaped {
				return "", ErrInvalidString
			}
		}

		if i+1 < len(src) && unicode.IsDigit(src[i+1]) {
			count, _ = strconv.Atoi(string(src[i+1]))
			i++
		}

		dst.WriteString(strings.Repeat(string(symbol), count))
	}

	return dst.String(), nil
}

// RuneEscaped Вернёт true если руна, в позиции i, из слайса r, экранирована
func RuneEscaped(r []rune, i int) (bool, error) {
	if i == 0 {
		return false, nil
	}
	if i < 0 || i >= len(r) {
		return false, ErrorOutOfRange
	}
	if r[i-1] != backslash {
		return false, nil
	}
	i--
	backSlashEscaped, _ := RuneEscaped(r, i)
	return !backSlashEscaped, nil
}
