package models

type GroupDBModel struct {
	Group   string   `db:"group_name"`
	Members []string `db:"members"`
}

func NewGroupDBModel(group string, members []string) GroupDBModel {
	return GroupDBModel{
		Group:   group,
		Members: members,
	}
}

type SongDBModel struct {
	ID         int64  `db:"id"`
	Song       string `db:"song"`
	Text       string `db:"text"`
	ReleasedAt string `db:"released_at"`
	Link       string `db:"link"`
}

func NewSongDBModel(song string, mod modelsHTTP) SongDBModel {
	return SongDBModel{
		ID:         mod.ID(),
		Song:       song,
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

func NewVerseDBModel(req VerseRequest) VerseDBModel {
	return VerseDBModel{
		ID:           req.ID,
		NumberVerse:  req.RequestedVerse,
		AmountVerses: 0,
		Text:         "",
	}
}

type RequestSongsDBModel struct {
	Limit  int64  `db:"limit"`
	Offset int64  `db:"offset"`
	Group  string `db:"group_name"`
	Song   string `db:"song"`
}

func NewSongsFilterDBModel(req RequestSongsFilter) RequestSongsDBModel {
	return RequestSongsDBModel{
		Limit:  req.Limit,
		Offset: req.Offset,
		Group:  req.Group,
		Song:   req.Song,
	}
}

type ResponceSongsDBModel struct {
	ID          int64  `db:"id"`
	Song        string `db:"song"`
	Group       string `db:"group_name"`
	ReleaseDate string `db:"released_at"`
	Link        string `db:"link"`
	Offset      int64  `db:"offset"`
}
