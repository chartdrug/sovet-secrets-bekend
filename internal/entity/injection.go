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
	What string `json:"what"`
}

type InjectionModel struct {
	Injection      Injection        `json:"injection"`
	Injection_Dose []Injection_Dose `json:"injection_dose"`
}

type PointValue struct {
	Drug string  `json:"drug"`
	C    float32 `json:"C"`
	CC   float32 `json:"CC"`
	CCT  float32 `json:"CCT"`
	CT   float32 `json:"CT"`
}

type Point struct {
	Dt          int          `json:"dt"`
	PointValues []PointValue `json:"pointValues"`
}
type PointsArray struct {
	Drugs  []string `json:"drugs"`
	Points []Point  `json:"point"`
}

// GetID returns the user ID.
func (u Injection) GetID() string {
	return u.ID
}
