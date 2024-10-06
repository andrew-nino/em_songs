package entity

type Group struct {
	Name    string
	Members []string	// TODO: add members table and struct
	Founded string
}

type Song struct {
	Name       string
	Text       string
	Album      string
	ReleasedAt string
	Link       string
}
