package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("read dir", func(t *testing.T) {
		var expected = Environment{
			"BAR":   "bar",
			"UNSET": "",
			"EMPTY": "",
			"FOO":   "   foo\nwith new line",
			"HELLO": "\"hello\"",
		}
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("path not found", func(t *testing.T) {
		var expected Environment
		env, err := ReadDir("./testdata/_env")
		require.Error(t, err)
		require.Equal(t, expected, env)
	})
}
