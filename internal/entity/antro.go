package entity

import "time"

//Login          string    `json:"login"`
type Antro struct {
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

	Fold_anterrior_iliac float32 `json:"fold_anterrior_iliac"`
	Fold_back            float32 `json:"fold_back"`
	Fold_belly           float32 `json:"fold_belly"`
	Fold_chest           float32 `json:"fold_chest"`
	Fold_forearm         float32 `json:"fold_forearm"`
	Fold_hip_front       float32 `json:"fold_hip_front"`
	Fold_hip_inside      float32 `json:"fold_hip_inside"`
	Fold_hip_rear        float32 `json:"fold_hip_rear"`
	Fold_hip_side        float32 `json:"fold_hip_side"`
	Fold_scapula         float32 `json:"fold_scapula"`
	Fold_shin            float32 `json:"fold_shin"`
	Fold_shoulder_front  float32 `json:"fold_shoulder_front"`
	Fold_shoulder_rear   float32 `json:"fold_shoulder_rear"`
	Fold_waist_side      float32 `json:"fold_waist_side"`
	Fold_wrist           float32 `json:"fold_wrist"`
	Fold_xiphoid         float32 `json:"fold_xiphoid"`

	Notes         string  `json:"notes"`
	Basic         bool    `json:"basic"`
	Result_fat    float32 `json:"result_fat"`
	Result_nofat  float32 `json:"result_nofat"`
	Result_energy float32 `json:"result_energy"`
}

// GetID returns the user ID.
func (u Antro) GetID() string {
	return u.ID
}
