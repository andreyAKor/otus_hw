package rmq

import (
	"time"
)

type Notification struct {
	EventID int64     `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UserID  int64     `json:"user_id"`
}
