package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const envDir = "./testdata/env"

func TestReadDir(t *testing.T) {
	envDir := "./testdata/env"

	t.Run("empty dir", func(t *testing.T) {
		_, err := ReadDir("")

		require.Error(t, err)
	})

	t.Run("file with \\n", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile(envDir, "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte("\n")); err != nil {
			t.Fatal(err)
		}

		fi, err := tmpfile.Stat()
		if err != nil {
			t.Fatal(err)
		}

		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		env, err := ReadDir(envDir)
		require.NoError(t, err)

		v, ok := env[fi.Name()]
		require.True(t, ok)
		require.Equal(t, "", v)
	})

	t.Run("file with zero byte", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile(envDir, "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte{0}); err != nil {
			t.Fatal(err)
		}

		fi, err := tmpfile.Stat()
		if err != nil {
			t.Fatal(err)
		}

		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		env, err := ReadDir(envDir)
		require.NoError(t, err)

		v, ok := env[fi.Name()]
		require.True(t, ok)
		require.Equal(t, "\n", v)
	})

	t.Run("filename with =", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile(envDir, "test=")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte("lala")); err != nil {
			t.Fatal(err)
		}

		fi, err := tmpfile.Stat()
		if err != nil {
			t.Fatal(err)
		}

		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		env, err := ReadDir(envDir)
		require.NoError(t, err)

		_, ok := env[fi.Name()]
		require.False(t, ok)
	})
}
