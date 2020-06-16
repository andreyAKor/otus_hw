package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	timeout time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout to host")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatal("programm syntax: go-telnet <host> <port>")
	}

	address := net.JoinHostPort(flag.Args()[0], flag.Args()[1])
	cl := NewTelnetClient(address, timeout, ioutil.NopCloser(os.Stdin), os.Stdout)

	log.Printf("...Connected to %s\n", address)

	if err := cl.Connect(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := cl.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			if err := cl.Send(); err != nil {
				log.Fatal(err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			if err := cl.Receive(); err != nil {
				log.Fatal(err)
				return
			}
		}
	}()

	wg.Wait()
}
