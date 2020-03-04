package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	tm, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatal("NTP error: ", err)
	}

	fmt.Println("current time:", time.Now())
	fmt.Println("exact time:", tm.UTC())
}
