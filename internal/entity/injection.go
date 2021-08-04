package entity

import "time"

//Login          string    `json:"login"`
type Injection struct {
	ID string

	Owner string    `json:"owner"`
	Dt    time.Time `json:"dt"`

	General_age       int     `json:"general_age"`
	General_hip       float32 `json:"general_hip"`
	General_height    float32 `json:"general_height"`
	General_leglen    float32 `json:"general_leglen"`
	General_weight    float32 `json:"general_weight"`
	General_handlen   float32 `json:"general_handlen"`
	General_shoulders float32 `json:"general_shoulders"`
	Notes             string  `json:"notes"`
	Basic             bool    `json:"basic"`
	Result_fat        float32 `json:"result_fat"`
	Result_nofat      float32 `json:"result_nofat"`
	Result_energy     float32 `json:"result_energy"`
}

// GetID returns the user ID.
func (u Injection) GetID() string {
	return u.ID
}
