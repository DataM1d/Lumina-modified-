package models

import "time"

type Analysis struct {
	ID          int64     `json:"id" db:"id"`
	URL         string    `json:"url" db:"url"`
	Headline    string    `json:"headline" db:"headline"`
	Summary     string    `json:"summary" db:"summary"`
	Sentiment   float64   `json:"sentiment" db:"sentiment"`
	VisualStyle string    `json:"visual_style" db:"visual_style"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

type Article struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	Analysis  Analysis  `json:"analysis"`
	CreatedAt time.Time `json:"created_at"`
}
