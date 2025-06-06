package model

import "time"

// UserEvent represents a Kinesis event
type UserEvent struct {
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}
