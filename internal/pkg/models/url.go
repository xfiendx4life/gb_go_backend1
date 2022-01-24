package models

type Url struct {
	Id           int       `json:"id"`
	Raw          string    `json:"raw"`
	Shortened    string    `json:"shortened"`
	UserId       int       `json:"user_id"`
	RedirectsNum Redirects `json:"redirects"`
}
