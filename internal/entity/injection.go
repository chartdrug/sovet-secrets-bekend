package entity

import "time"

//Login          string    `json:"login"`

type Injection_Dose struct {
	ID           string  `json:"id"`
	Id_injection string  `json:"-"`
	Dose         float64 `json:"dose"`
	Drug         string  `json:"drug"`
	Volume       float64 `json:"volume"`
	Solvent      string  `json:"solvent"`
}

type Injection struct {
	ID       string
	Owner    string    `json:"-"`
	Dt       time.Time `json:"dt"`
	Course   string    `json:"course"`
	What     string    `json:"what"`
	TotalV   float64   `json:"-" db:"-"`
	SkinSumm float64   `json:"-" db:"-"`
	Skin     float64   `json:"-" db:"-"`
	Calc     bool      `json:"calc"`
}

type Injectionv2 struct {
	ID       string
	Owner    string   `json:"-"`
	Course   string   `json:"course"`
	Date     []string `json:"date"`
	Times    []int    `json:"times"`
	What     string   `json:"what"`
	TotalV   float64  `json:"-" db:"-"`
	SkinSumm float64  `json:"-" db:"-"`
	Skin     float64  `json:"-" db:"-"`
	Calc     bool     `json:"-"`
}

type Concentration struct {
	Owner        string
	Id_injection string
	Drug         string
	Dt           int64
	C            float64
	CC           float64
	CCT          float64
	CT           float64
}

type InjectionModel struct {
	Injection      Injection        `json:"injection"`
	Injection_Dose []Injection_Dose `json:"injection_dose"`
}

type InjectionModelv2 struct {
	Injection      Injectionv2      `json:"injection"`
	Injection_Dose []Injection_Dose `json:"injection_dose"`
}

type PointValue struct {
	Drug string  `json:"drug"`
	C    float64 `json:"-"`
	CC   float64 `json:"-"`
	CCT  float64 `json:"CCT"`
	CT   float64 `json:"-"`
	R    float64 `json:"-"`
	//CO  float64
	//COT  float64
	OutK   float64 `json:"-" db:"-"`
	OutKT  float64 `json:"-" db:"-"`
	Dose   float64 `json:"-"`
	Volume float64 `json:"-"`
	Ri     float64 `json:"-"`
	Depo   float64 `json:"-"`
	Depoi  float64 `json:"-"`
	Dv     float64 `json:"-"`
	Cout   float64 `json:"-"`
	Coutt  float64 `json:"-"`
}

type Point struct {
	Dt          int64        `json:"dt"`
	PointValues []PointValue `json:"pointValues"`
}
type PointsArray struct {
	CountProcess int              `json:"countProcess"`
	Pkt          time.Time        `json:"pkt"`
	Control      time.Time        `json:"control"`
	Drugs        []string         `json:"drugs"`
	Points       []Point          `json:"point"`
	Injections   []InjectionModel `json:"Injections"`
}

// GetID returns the user ID.
func (u Injection) GetID() string {
	return u.ID
}

/*
type concentration struct {
	Drug string
}*/
