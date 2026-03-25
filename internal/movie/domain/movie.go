package domain

import "github.com/google/uuid"

type Movie struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	ReleaseYear  int       `json:"release_year"`
	Rating       float64   `json:"rating"`
	Duration     int       `json:"duration"`
	BriefContent string    `json:"brief_content"`
	PreviewURL   string    `json:"preview_url"`
	MovieURL     string    `json:"movie_url"`
	CoverURL     string    `json:"cover_url"`
}
