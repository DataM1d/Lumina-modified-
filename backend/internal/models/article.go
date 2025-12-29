package models

import "time"

type Article struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Content     string    `json:"content"`
	Headline    string    `json:"headline"`
	Summary     string    `json:"summary"`
	Sentiment   float64   `json:"sentiment"`
	VisualStyle string    `json:"visual_style"`
	CreatedAt   time.Time `json:"created_at"`
}
