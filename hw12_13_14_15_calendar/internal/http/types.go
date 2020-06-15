package http

import (
	"time"
)

// Response is the struct of response.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// Event is the struct of event.
type Event struct {
	Title         string    `json:"title"`
	Date          time.Time `json:"date"`
	Duration      string    `json:"duration"`
	Descr         *string   `json:"descr,omitempty"`
	UserID        int64     `json:"user_id"`
	DurationStart *string   `json:"duration_start,omitempty"`
}

type EventList struct {
	Events []Event `json:"events"`
}
