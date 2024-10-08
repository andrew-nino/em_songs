package service

import (
	"context"

	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/sirupsen/logrus"
)

type SongsRepo interface {
	AddSongToRepository(context.Context, models.GroupDBModel, models.SongDBModel) (int, error)
	UpdateSongToRepository(context.Context, models.SongDBModel) error
	DeleteSongFromRepository(context.Context, int) error
}

type ApplicationServices struct {
	log        *logrus.Logger
	repository SongsRepo
}

func New(log *logrus.Logger, songs SongsRepo) *ApplicationServices {
	return &ApplicationServices{
		log:        log,
		repository: songs,
	}
}
