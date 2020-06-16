package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("empty address", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("", timeout, ioutil.NopCloser(in), out)
		require.Error(t, client.Connect())
		defer func() {
			require.NoError(t, client.Close())
		}()

		in.WriteString("hello\n")
		err = client.Send()
		require.NoError(t, err)

		err = client.Receive()
		require.NoError(t, err)
	})

	t.Run("io timeout 1 sec for IP", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("1s")
		require.NoError(t, err)

		client := NewTelnetClient("128.0.0.1:80", timeout, ioutil.NopCloser(in), out)

		start := time.Now()
		err = client.Connect()
		elapsedTime := time.Since(start)

		require.Error(t, err)
		require.LessOrEqual(t, int64(elapsedTime), int64(float64(timeout)*1.1))

		require.NoError(t, client.Close())
	})

	t.Run("io timeout 1 sec for domain", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("1s")
		require.NoError(t, err)

		client := NewTelnetClient("lala.ru:80", timeout, ioutil.NopCloser(in), out)

		start := time.Now()
		err = client.Connect()
		elapsedTime := time.Since(start)

		require.Error(t, err)
		require.LessOrEqual(t, int64(elapsedTime), int64(float64(timeout)*1.1))

		require.NoError(t, client.Close())
	})

	t.Run("connection refused", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("1s")
		require.NoError(t, err)

		client := NewTelnetClient("127.0.0.1:65535", timeout, ioutil.NopCloser(in), out)
		require.Error(t, client.Connect())
		require.NoError(t, client.Close())
	})

	t.Run("server shutdown", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.EqualError(t, ErrConnClose, err.Error())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)

			require.NoError(t, conn.Close())
		}()

		wg.Wait()
	})

	t.Run("client shutdown", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			_, err = conn.Read(request)
			require.NoError(t, err)

			_, err = conn.Read(request)
			require.EqualError(t, err, io.EOF.Error())
		}()

		wg.Wait()
	})
}
