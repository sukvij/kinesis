package model

import "time"

// User represents a user in the database
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
