package entity

import "time"

//Login          string    `json:"login"`
type Antro struct {
	ID string

	Owner string    `json:"owner"`
	Dt    time.Time `json:"dt"`

	General_age       int     `json:"general_age"`
	General_hip       float64 `json:"general_hip"`
	General_height    float64 `json:"general_height"`
	General_leglen    float64 `json:"general_leglen"`
	General_weight    float64 `json:"general_weight"`
	General_handlen   float64 `json:"general_handlen"`
	General_shoulders float64 `json:"general_shoulders"`

	Fold_anterrior_iliac float64 `json:"fold_anterrior_iliac"`
	Fold_back            float64 `json:"fold_back"`
	Fold_belly           float64 `json:"fold_belly"`
	Fold_chest           float64 `json:"fold_chest"`
	Fold_forearm         float64 `json:"fold_forearm"`
	Fold_hip_front       float64 `json:"fold_hip_front"`
	Fold_hip_inside      float64 `json:"fold_hip_inside"`
	Fold_hip_rear        float64 `json:"fold_hip_rear"`
	Fold_hip_side        float64 `json:"fold_hip_side"`
	Fold_scapula         float64 `json:"fold_scapula"`
	Fold_shin            float64 `json:"fold_shin"`
	Fold_shoulder_front  float64 `json:"fold_shoulder_front"`
	Fold_shoulder_rear   float64 `json:"fold_shoulder_rear"`
	Fold_waist_side      float64 `json:"fold_waist_side"`
	Fold_wrist           float64 `json:"fold_wrist"`
	Fold_xiphoid         float64 `json:"fold_xiphoid"`

	Notes         string  `json:"notes"`
	Basic         bool    `json:"basic"`
	Result_fat    float64 `json:"result_fat"`
	Result_nofat  float64 `json:"result_nofat"`
	Result_energy float64 `json:"result_energy"`
}

// GetID returns the user ID.
func (u Antro) GetID() string {
	return u.ID
}
