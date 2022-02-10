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
	TypeSports     []int     `json:"typesports"`
	Admin          bool      `json:"admin"`
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
	Admin          bool
}

type UsersAdmin struct {
	Id              string
	Login           string
	Admin           bool
	Blocked         bool
	Blocked_at      time.Time
	Email           string
	Date_registered time.Time
	Date_lastlogin  time.Time
	Sex             string
	Birthday        time.Time
	Type_sports     string
	Country         string
	Region          string
	City            string
	Count_antro     int
	Count_injection int
}

type HistoryLogin struct {
	IdUser    string
	DateEvent time.Time `json:"DateEvent"`
	IpAddress string    `json:"IpAddress"`
	Country   string    `json:"Country"`
	Region    string    `json:"Region"`
	City      string    `json:"City"`
}

type Feedback struct {
	Id       string    `json:"id"`
	Owner    string    `json:"owner"`
	Dt       time.Time `json:"dt"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Feedback string    `json:"feedback"`
	Location string    `json:"location"`
}

// GetID returns the user ID.
func (u Users) GetID() string {
	return u.ID
}
func (u Users) GetLogn() string {
	return u.Login
}

func (u Users) GetAdmin() bool {
	return u.Admin
}

func (u Users) GetSex() string {
	return u.Sex
}
