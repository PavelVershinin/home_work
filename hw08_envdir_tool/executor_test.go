package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("remove empty", func(t *testing.T) {
		stdOut := bytes.Buffer{}
		env := Environment{
			"FOO": "",
		}

		require.NoError(t, os.Setenv("FOO", "BAR"))
		ret := runCmd("env", nil, env, nil, &stdOut, nil)

		require.Equal(t, ret, 0)
		require.NotContains(t, stdOut.String(), "FOO")
	})

	t.Run("override", func(t *testing.T) {
		stdOut := bytes.Buffer{}
		env := Environment{
			"FOO": "OVERRIDDEN",
		}

		require.NoError(t, os.Setenv("FOO", "BAR"))
		ret := runCmd("env", nil, env, nil, &stdOut, nil)

		require.Equal(t, ret, 0)
		require.Contains(t, stdOut.String(), "FOO=OVERRIDDEN")
	})
}
