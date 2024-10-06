package service

import (
	"context"
	"encoding/json"

	"github.com/andrew-nino/em_songs/internal/models"
)

func (s *ApplicationServices) AddSong(ctx context.Context, songRequest models.SongRequest, rawAnswer []byte) (int, error) {

	requestModel, err := processRawAnswer(ctx, rawAnswer)
	if err != nil {
		s.log.WithError(err).Error("Failed to process raw answer")
		return 0, err
	}

	groupDBModel := models.NewGroupDBModel(songRequest.Group, nil)
	songModel := models.NewSongDBModel(songRequest.Song, requestModel)

	id, err := s.songs.AddSongToRepository(ctx, *groupDBModel, *songModel)
	if err != nil {
		s.log.WithError(err).Error("Failed to add song to repository")
		return 0, err
	}

	return id, nil
}
func processRawAnswer(ctx context.Context, body []byte) (models.SongDetail, error) {

	var songDetail models.SongDetail

	done := make(chan error, 1)

	go func() {
		defer close(done)

		err := json.Unmarshal(body, &songDetail)
		if err != nil {
			done <- err
			return
		}
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return models.SongDetail{}, err
		}
		return songDetail, nil
	case <-ctx.Done():
		return models.SongDetail{}, ctx.Err()
	}
}
