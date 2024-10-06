package models

type SongDetail struct {
	ReleaseDate string `json:"releaseDate" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link" binding:"required"`
}

type SongRequest struct {
	Group string `json:"group" validate:"required,lte=100"`
	Song  string `json:"song" validate:"required,lte=100"`
}
