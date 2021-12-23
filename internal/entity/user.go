package entity

import "time"

// User represents a user.
type UsersRest struct {
	ID             string
	Login          string `json:"login"`
	Passwd         string
	Email          string    `json:"email"`
	DateRegistered time.Time `json:"date_registered"`
	DateLastlogin  time.Time `json:"date_lastlogin"`
	Sex            string    `json:"sex"`
	Birthday       time.Time `json:"birthday"`
	TypeSports     []string  `json:"typesports"`
}

type Users struct {
	ID             string
	Login          string `json:"login"`
	Passwd         string
	Email          string    `json:"email"`
	DateRegistered time.Time `json:"date_registered"`
	DateLastlogin  time.Time `json:"date_lastlogin"`
	Sex            string    `json:"sex"`
	Birthday       time.Time `json:"birthday"`
	TypeSports     string    `json:"typesports"`
}

type HistoryLogin struct {
	IdUser    string
	DateEvent time.Time `json:"DateEvent"`
	IpAddress string    `json:"IpAddress"`
	Country   string    `json:"Country"`
	Region    string    `json:"Region"`
	City      string    `json:"City"`
}

// GetID returns the user ID.
func (u Users) GetID() string {
	return u.ID
}
func (u Users) GetLogn() string {
	return u.Login
}
