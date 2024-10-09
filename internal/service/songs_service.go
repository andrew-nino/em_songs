package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/andrew-nino/em_songs/internal/models"
)

func (s *ApplicationServices) AddSong(ctx context.Context, songRequest models.SongRequest, rawAnswer []byte) (int, error) {

	var songDetail models.SongDetail
	err := processRawAnswer(&songDetail, rawAnswer)
	if err != nil {
		s.log.WithError(err).Error("failed to process raw answer")
		return 0, err
	}

	groupDBModel := models.NewGroupDBModel(songRequest.Group, nil)
	songModel := models.NewSongDBModel(songRequest.Song, songDetail)

	id, err := s.repository.AddSongToRepository(ctx, *groupDBModel, *songModel)
	if err != nil {
		s.log.WithError(err).Error("failed to add song to repository")
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

func (s *ApplicationServices) UpdateSong(ctx context.Context, updateSong models.SongUpdate) error {

	songModel := models.NewSongDBModel(updateSong.Name, updateSong)

	err := s.repository.UpdateSongToRepository(ctx, *songModel)
	if err != nil {
		s.log.WithError(err).Error("failed to update song in repository")
		return err
	}
	return nil
}

func (s *ApplicationServices) GetSong(ctx context.Context, request models.VerseRequest) (models.VerseResponce, error) {

	verseDBModel := models.NewVerseDBModel(request)
	responceFromDB, err := s.repository.GetSong(ctx, *verseDBModel)
	if err != nil {
		s.log.WithError(err).Error("failed to get song from repository")
		return models.VerseResponce{}, err
	}

	sliceVerses := strings.Split(responceFromDB.Text, "\n\n")
	lenSliceVerses := len(sliceVerses)

	if lenSliceVerses <= int(request.RequestedVerse) || request.RequestedVerse <= zero {
		request.RequestedVerse = one
	}

	var nextVesre int64

	if lenSliceVerses == zero {
		return models.VerseResponce{}, fmt.Errorf("the requested song does not have any verses")
	} else if lenSliceVerses > int(request.RequestedVerse + one) {
		nextVesre = request.RequestedVerse + one
	} else if lenSliceVerses == int(request.RequestedVerse) {
		nextVesre = zero
	}

	responce := models.VerseResponce{
		NextVerse: nextVesre,
		Text:      sliceVerses[request.RequestedVerse - one],
	}

	return responce, nil
}

func (s *ApplicationServices) DeleteSong(ctx context.Context, id int) error {

	err := s.repository.DeleteSongFromRepository(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("failed to delete song from repository")
		return err
	}
	return nil
}
