package postgres

import (
	"context"

	"github.com/andrew-nino/em_songs/internal/models"
)

func (p *Postgres) AddSongToRepository(context.Context, models.GroupDBModel, models.SongDBModel) (int, error) {

	return 0, nil
}
