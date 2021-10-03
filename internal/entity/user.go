package entity

import "time"

// User represents a user.
type Users struct {
	ID             string
	Login          string `json:"login"`
	Passwd         string
	Email          string    `json:"email"`
	DateRegistered time.Time `json:"date_registered"`
	DateLastlogin  time.Time
	Sex            string    `json:"sex"`
	Birthday       time.Time `json:"birthday"`
	Height         float64   `json:"height"`
	Weight         float64   `json:"weight"`
}

// GetID returns the user ID.
func (u Users) GetID() string {
	return u.ID
}
func (u Users) GetLogn() string {
	return u.Login
}
