package hw02_unpack_string //nolint:golint,stylecheck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abcd",
			expected: "abcd",
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestRuneEscaped(t *testing.T) {
	var testTable = []struct {
		input    []rune
		position int
		expected bool
		err      error
	}{
		{
			input:    []rune(`Test 符号 1 \`),
			position: 0,
			expected: false,
			err:      nil,
		},
		{
			input:    []rune(`Test 符号 1 \`),
			position: -1,
			expected: false,
			err:      ErrorOutOfRange,
		},
		{
			input:    []rune(`Test 符号 1 \`),
			position: 10,
			expected: false,
			err:      nil,
		},
		{
			input:    []rune(`Test 符号 1 \`),
			position: 11,
			expected: false,
			err:      ErrorOutOfRange,
		},
		{
			input:    []rune(`Test 符\号 1 \`),
			position: 7,
			expected: true,
			err:      nil,
		},
		{
			input:    []rune(`Test 符\\号 1 \`),
			position: 8,
			expected: false,
			err:      nil,
		},
		{
			input:    []rune(`Test 符\\号 1 \`),
			position: 7,
			expected: true,
			err:      nil,
		},
	}

	for i, tst := range testTable {
		t.Run(fmt.Sprintf("test #%d", i), func(t *testing.T) {
			result, err := RuneEscaped(tst.input, tst.position)
			require.Equal(t, tst.err, err)
			require.Equal(t, tst.expected, result)
		})
	}
}
