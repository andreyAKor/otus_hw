package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("return 0", func(t *testing.T) {
		r := RunCmd([]string{"env", "-i"}, Environment{})
		require.Equal(t, 0, r)
	})

	t.Run("run unknow command", func(t *testing.T) {
		r := RunCmd([]string{"abracadabra"}, Environment{})
		require.Equal(t, -1, r)
	})

	t.Run("run command without params", func(t *testing.T) {
		r := RunCmd([]string{"cp"}, Environment{})
		require.Equal(t, 1, r)
	})
}
