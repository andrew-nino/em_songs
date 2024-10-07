package service

import (
	"context"
	"encoding/json"

	"github.com/andrew-nino/em_songs/internal/models"
)

func (s *ApplicationServices) AddSong(ctx context.Context, songRequest models.SongRequest, rawAnswer []byte) (int, error) {

	var songDetail models.SongDetail
	err := processRawAnswer(&songDetail, rawAnswer)
	if err != nil {
		s.log.WithError(err).Error("Failed to process raw answer")
		return 0, err
	}

	groupDBModel := models.NewGroupDBModel(songRequest.Group, nil)
	songModel := models.NewSongDBModel(songRequest.Song, songDetail)

	id, err := s.songs.AddSongToRepository(ctx, *groupDBModel, *songModel)
	if err != nil {
		s.log.WithError(err).Error("Failed to add song to repository")
		return 0, err
	}

	return id, nil
}

func processRawAnswer(songDetail *models.SongDetail, body []byte) error {

	err := json.Unmarshal(body, &songDetail)
	if err != nil {
		return err
	}
	return nil
}
