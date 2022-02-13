package models

import "time"

type Summary struct {
	Month int `json:"month"`
	Week  int `json:"week"`
	Today int `json:"today"`
	// Add Last Usage Date
}

type Redirects struct {
	Id    int
	UrlId int
	Date  time.Time
}
