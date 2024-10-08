package service

import (
	"context"

	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/sirupsen/logrus"
)

type SongsRepo interface {
	AddSongToRepository(context.Context, models.GroupDBModel, models.SongDBModel) (int, error)
	DeleteSongFromRepository(context.Context, int) error
}

type ApplicationServices struct {
	log   *logrus.Logger
	songs SongsRepo
}

func New(log *logrus.Logger, songs SongsRepo) *ApplicationServices {
	return &ApplicationServices{
		log:   log,
		songs: songs,
	}
}
