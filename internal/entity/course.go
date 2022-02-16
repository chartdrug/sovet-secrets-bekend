package entity

import "time"

// User represents a user.
type Course struct {
	Id           string    `json:"id"`
	Owner        string    `json:"owner"`
	Descr        string    `json:"descr"`
	Course_start time.Time `json:"course_start"`
	Course_end   time.Time `json:"course_end"`
	Type         string    `json:"type"`
	Target       []string  `json:"target"`
	Notes        string    `json:"notes"`
}

// GetID returns the user ID.
func (u Course) GetID() string {
	return u.Id
}
