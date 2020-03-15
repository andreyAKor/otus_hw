package hw02_unpack_string //nolint:golint,stylecheck

import (
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
			input:    "4a5b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "a",
			expected: "a",
		},
		{
			input:    "a9",
			expected: "aaaaaaaaa",
		},
		{
			input:    "a0b1",
			expected: "ab",
		},
		{
			input:    "a-1b-2c+3",
			expected: "a-b--c+++",
		},
		{
			input:    "*@0&1%2♬♬3",
			expected: "*@&%%♬♬♬♬",
		},
		{
			input:    string('0'),
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    string(0x0),
			expected: string(0x0),
		},
		{
			input:    string(0x0 - 1),
			expected: string(0x0 - 1),
		},
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
	t.Skip() // Remove if task with asterisk completed

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
