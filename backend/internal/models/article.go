package models

import "time"

type Analysis struct {
	Headline    string    `json:"headline"`
	Summary     string    `json:"summary"`
	Sentiment   float64   `json:"sentiment"`
	VisualStyle string    `json:"visual_style"`
	ProcessedAt time.Time `json:"processed_at"`
}

type Article struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	Analysis  Analysis  `json:"analysis"`
	CreatedAt time.Time `json:"created_at"`
}
