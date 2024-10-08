package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/andrew-nino/em_songs/internal/models"
)

func (p *Postgres) AddSongToRepository(ctx context.Context, group models.GroupDBModel, song models.SongDBModel) (int, error) {

	var id_group int
	var id_song int
	var operationID int

	// Проверяем наличие группы в таблице
	searchGroupQuery := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", groups_table)
	rowGroup := p.db.QueryRowContext(ctx, searchGroupQuery, group.Name)
	err := rowGroup.Scan(&id_group)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	if id_group == 0 {

		insertGroupQuery := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", groups_table)
		row := tx.QueryRowContext(ctx, insertGroupQuery, group.Name)
		if err = row.Scan(&id_group); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	insertSongQuery := fmt.Sprintf("INSERT INTO %s (name, text, released_at, link) VALUES ($1, $2, $3, $4) RETURNING id", songs_table)
	rowSong := tx.QueryRowContext(ctx, insertSongQuery, song.Name, song.Text, song.ReleasedAt, song.Link)
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

func (p *Postgres) DeleteSongFromRepository(ctx context.Context, id int) error {

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", songs_table)
	_, err := p.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	return nil
}
