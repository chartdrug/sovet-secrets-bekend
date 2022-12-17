package entity

import (
	"time"
)

// Album represents an album record.
type HistoryUpdate struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	DescriptionRu string    `json:"description_ru"`
	DescriptionEn string    `json:"description_en"`
}
