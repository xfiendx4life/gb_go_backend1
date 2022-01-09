package models

type Redirects struct {
	Month int `json:"month"`
	Week  int `json:"week"`
	Today int `json:"today"`
	// Add Last Usage Date
}
