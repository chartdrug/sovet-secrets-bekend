package entity

import "time"

//Login          string    `json:"login"`

type Injection_Dose struct {
	ID           string  `json:"id"`
	Id_injection string  `json:"id_injection"`
	Dose         float32 `json:"dose"`
	Drug         string  `json:"drug"`
	Volume       float32 `json:"volume"`
	Solvent      string  `json:"solvent"`
}

type Injection struct {
	ID    string
	Owner string    `json:"owner"`
	Dt    time.Time `json:"dt"`
	//Course  string           `json:"course"`
	What           string           `json:"what"`
	Injection_Dose []Injection_Dose `json:"injection_dose"`
}

// GetID returns the user ID.
func (u Injection) GetID() string {
	return u.ID
}
