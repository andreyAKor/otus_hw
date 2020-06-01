package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrConnClose = errors.New("...Connection was closed by peer")
	ErrEOF       = errors.New("...EOF")
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

func NewTelnetClient(
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type Telnet struct {
	address     string
	timeout     time.Duration
	in          io.ReadCloser
	out         io.Writer
	conn        net.Conn
	connScanner *bufio.Scanner
	inScanner   *bufio.Scanner
}

func (t *Telnet) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)

	t.connScanner = bufio.NewScanner(t.conn)
	t.inScanner = bufio.NewScanner(t.in)

	return
}

func (t *Telnet) Receive() (err error) {
	if t.conn == nil {
		return
	}
	if !t.connScanner.Scan() {
		return ErrConnClose
	}
	_, err = t.out.Write([]byte(t.connScanner.Text() + "\n"))
	return
}

func (t *Telnet) Send() (err error) {
	if t.conn == nil {
		return
	}
	if !t.inScanner.Scan() {
		return ErrEOF
	}
	_, err = t.conn.Write([]byte(t.inScanner.Text() + "\n"))
	return
}

func (t *Telnet) Close() (err error) {
	if t.conn != nil {
		err = t.conn.Close()
	}
	return
}

var _ TelnetClient = (*Telnet)(nil)
