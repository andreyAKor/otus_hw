package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const delay = 5 * time.Second

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	// Calendar
	calendarStatus := godog.RunWithOptions("calendar integration", CalendarFeatureContext, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features/calendar.feature"},
		Randomize: 0,
	})
	if calendarStatus > 0 {
		os.Exit(calendarStatus)
		return
	}

	// Sender
	senderStatus := godog.RunWithOptions("sender integration", SenderFeatureContext, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features/sender.feature"},
		Randomize: 0,
	})
	if senderStatus > 0 {
		os.Exit(senderStatus)
		return
	}

	if st := m.Run(); st > 0 {
		os.Exit(st)
		return
	}

	os.Exit(0)
}
