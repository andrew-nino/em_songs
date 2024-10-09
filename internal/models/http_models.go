package models

type SongRequest struct {
	Group string `json:"group" validate:"required,lte=100"`
	Song  string `json:"song" validate:"required,lte=100"`
}

type modelsHTTP interface {
	ID() int64
	Text() string
	ReleaseDate() string
	Reference() string
}

type SongDetail struct {
	SongID   int64  `json:"id"`
	Release  string `json:"releaseDate" binding:"required"`
	TextSong string `json:"text" binding:"required"`
	Link     string `json:"link" binding:"required"`
}

func (sd SongDetail) ID() int64 {
	return sd.SongID
}
func (sd SongDetail) Text() string {
	return sd.TextSong
}
func (sd SongDetail) ReleaseDate() string {
	return sd.Release
}
func (sd SongDetail) Reference() string {
	return sd.Link
}

type SongUpdate struct {
	SongID   int64  `json:"id" validate:"required"`
	Name     string `json:"name" validate:"lte=100"`
	Release  string `json:"releaseDate"`
	TextSong string `json:"text"`
	Link     string `json:"link"`
}

func (su SongUpdate) ID() int64 {
	return su.SongID
}
func (su SongUpdate) Text() string {
	return su.TextSong
}
func (su SongUpdate) ReleaseDate() string {
	return su.Release
}
func (su SongUpdate) Reference() string {
	return su.Link
}

type VerseRequest struct {
	ID             int64 `json:"id" validate:"required"`
	RequestedVerse int64 `json:"requestedVerse"`
}

type VerseResponce struct {
	NextVerse int64  `json:"nextVerse"`
	Text      string `json:"text"`
}
