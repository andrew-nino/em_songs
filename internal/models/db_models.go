package models

type GroupDBModel struct {
	Name    string   `db:"name"`
	Members []string `db:"members"`
}

func NewGroupDBModel(name string, members []string) *GroupDBModel {
	return &GroupDBModel{
		Name:    name,
		Members: members,
	}
}

type SongDBModel struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	Text       string `db:"text"`
	ReleasedAt string `db:"released_at"`
	Link       string `db:"link"`
}

func NewSongDBModel(name string, mod modelsHTTP) *SongDBModel {
	return &SongDBModel{
		ID:         mod.ID(),
		Name:       name,
		Text:       mod.Text(),
		ReleasedAt: mod.ReleaseDate(),
		Link:       mod.Reference(),
	}
}

type VerseDBModel struct {
	ID           int64  `db:"id"`
	NumberVerse  int64  // for Redis
	AmountVerses int64  // for Redis
	Text         string `db:"text"`
}

func NewVerseDBModel(req VerseRequest) *VerseDBModel {
	return &VerseDBModel{
		ID:           req.ID,
		NumberVerse:  req.RequestedVerse,
		AmountVerses: 0,
		Text:         "",
	}
}
