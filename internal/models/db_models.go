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
	Name       string `db:"name"`
	Text       string `db:"text"`
	ReleasedAt string `db:"released_at"`
	Link       string `db:"link"`
}

func NewSongDBModel(name string, mod SongDetail) *SongDBModel {
	return &SongDBModel{
		Name:       name,
		Text:       mod.Text,
		ReleasedAt: mod.ReleaseDate,
		Link:       mod.Link,
	}
}
