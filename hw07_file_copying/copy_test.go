package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("empty files", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		require.True(t, os.IsNotExist(err))

		err = Copy("/etc/hosts", "", 0, 0)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("unlimit file", func(t *testing.T) {
		err := Copy("/dev/urandom", "", 0, 0)

		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write([]byte("")); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		err = Copy(tmpfile.Name(), "", 1, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("permission file", func(t *testing.T) {
		err := Copy("/dev/psaux", "", 0, 0)

		require.True(t, os.IsPermission(err))
	})
}
