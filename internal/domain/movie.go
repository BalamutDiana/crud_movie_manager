package domain

import "time"

type Movie struct {
	ID               int       `json:"id,omitempty"`
	Title            string    `json:"title"`
	Release          string    `json:"release"`
	StreamingService string    `json:"streamingService"`
	SavedAt          time.Time `json:"savedAt,omitempty"`
}
