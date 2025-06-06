package model

import "time"

// Log represents a log entry in the database
type Log struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Action    string    `gorm:"not null"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}
