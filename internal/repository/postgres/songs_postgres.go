package postgres

import (
	"context"
	"database/sql"
	"fmt"

	v1 "github.com/andrew-nino/em_songs/internal/controller/http/v1"
	"github.com/andrew-nino/em_songs/internal/models"
)

func (p *Postgres) AddSongToRepository(ctx context.Context, gr models.GroupDBModel, sng models.SongDBModel) (int, error) {

	var id_group int
	var id_song int
	var operationID int

	// Проверяем наличие группы в таблице
	searchGroupQuery := fmt.Sprintf("SELECT id FROM %s WHERE group_name = $1", groups_table)
	rowGroup := p.db.QueryRowContext(ctx, searchGroupQuery, gr.Group)
	err := rowGroup.Scan(&id_group)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	if id_group == 0 {

		insertGroupQuery := fmt.Sprintf("INSERT INTO %s (group_name) VALUES ($1) RETURNING id", groups_table)
		row := tx.QueryRowContext(ctx, insertGroupQuery, gr.Group)
		if err = row.Scan(&id_group); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	insertSongQuery := fmt.Sprintf("INSERT INTO %s (song, text, released_at, link) VALUES ($1, $2, $3, $4) RETURNING id", songs_table)
	rowSong := tx.QueryRowContext(ctx, insertSongQuery, sng.Song, sng.Text, sng.ReleasedAt, sng.Link)
	if err = rowSong.Scan(&id_song); err != nil {
		tx.Rollback()
		return 0, err
	}

	insertGroupSongQuery := fmt.Sprintf("INSERT INTO %s (group_id, song_id) VALUES ($1, $2) RETURNING id", group_song_table)
	rowGroupSong := tx.QueryRowContext(ctx, insertGroupSongQuery, id_group, id_song)
	if err = rowGroupSong.Scan(&operationID); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id_song, tx.Commit()
}

func (p *Postgres) UpdateSongToRepository(ctx context.Context, songUpdate models.SongDBModel) error {

	var id_song int64

	_, err := p.db.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		return err
	}

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	insertSongQuery := fmt.Sprintf("UPDATE %s SET song=$1, text=$2, released_at=$3, link=$4, updated_at=now() WHERE id = $5 RETURNING id", songs_table)
	rowSong := tx.QueryRowContext(ctx, insertSongQuery, songUpdate.Song, songUpdate.Text, songUpdate.ReleasedAt, songUpdate.Link, songUpdate.ID)
	if err = rowSong.Scan(&id_song); err != nil && id_song != songUpdate.ID {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *Postgres) GetSong(ctx context.Context, verseDBModel models.VerseDBModel) (models.VerseDBModel, error) {

	var verses []string
	getSongQuery := fmt.Sprintf("SELECT text FROM %s WHERE id = $1", songs_table)
	err := p.db.SelectContext(ctx, &verses, getSongQuery, verseDBModel.ID)
	if err != nil {
		return models.VerseDBModel{}, err
	}

	if len(verses) == 0 {
		return models.VerseDBModel{}, v1.ErrNoSongFound
	} else {
		verseDBModel.Text = verses[0]
	}

	return verseDBModel, nil
}

func (p *Postgres) GetAllSongs(ctx context.Context, requestSongs models.RequestSongsDBModel) ([]models.ResponceSongsDBModel, error) {

	var songs = make([]models.ResponceSongsDBModel, 0)
	var query string

	if requestSongs.Group == "" {

		query = fmt.Sprintf(`SELECT s.id, s.song, g.group_name, s.released_at, s.link 
							 FROM %s AS s 
							 INNER JOIN %s AS gs ON s.id = gs.song_id	
							 INNER JOIN %s AS g ON gs.group_id = g.id 
							 WHERE s.song ILIKE '%%'||$1||'%%' AND s.id >= $2 
							 ORDER BY s.id LIMIT $3`, songs_table, group_song_table, groups_table)

		err := p.db.SelectContext(ctx, &songs, query, requestSongs.Song, requestSongs.Offset, requestSongs.Limit)
		if err != nil {
			return nil, err
		}

	} else {
		query = fmt.Sprintf(`SELECT s.id, s.song, g.group_name, s.released_at, s.link 
							 FROM %s AS s 
							 INNER JOIN %s AS gs ON s.id = gs.song_id	
							 INNER JOIN %s AS g ON gs.group_id = g.id AND g.group_name = $1
							 WHERE s.song ILIKE '%%'||$2||'%%' AND s.id >= $3 
							 ORDER BY s.id LIMIT $4`, songs_table, group_song_table, groups_table)

		err := p.db.SelectContext(ctx, &songs, query, requestSongs.Group, requestSongs.Song, requestSongs.Offset, requestSongs.Limit)
		if err != nil {
			return nil, err
		}
	}

	return songs, nil
}

func (p *Postgres) DeleteSongFromRepository(ctx context.Context, id int) error {

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", songs_table)
	_, err := p.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	return nil
}
