package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/sirupsen/logrus"
)

func (s *ApplicationServices) AddSong(ctx context.Context, songRequest models.SongRequest, rawAnswer []byte) (int, error) {
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "AddSong"}).Infof("song:%v", songRequest.Song)

	var songDetail models.SongDetail
	err := processRawAnswer(&songDetail, rawAnswer)
	if err != nil {
		s.log.WithError(err).Error("failed to process raw answer")
		return 0, err
	}

	groupDBModel := models.NewGroupDBModel(songRequest.Group, nil)
	songModel := models.NewSongDBModel(songRequest.Song, songDetail)

	id, err := s.repository.AddSongToRepository(ctx, groupDBModel, songModel)
	if err != nil {
		s.log.WithError(err).Error("failed to add song to repository")
		return 0, err
	}
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "AddSong"}).Infof("success. songID:%v", id)

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
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "UpdateSong"}).Infof("songID:%v", updateSong.SongID)

	songModel := models.NewSongDBModel(updateSong.Song, updateSong)
	err := s.repository.UpdateSongToRepository(ctx, songModel)
	if err != nil {
		s.log.WithError(err).Error("failed to update song in repository")
		return err
	}
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "UpdateSong"}).Infof("success. songID:%v", updateSong.SongID)

	return nil
}

func (s *ApplicationServices) GetSong(ctx context.Context, request models.VerseRequest) (models.VerseResponce, error) {
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "GetSong"}).Infof("request:%+v", request)

	verseDBModel := models.NewVerseDBModel(request)
	responceFromDB, err := s.repository.GetSong(ctx, verseDBModel)
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
	} else if lenSliceVerses > int(request.RequestedVerse+one) {
		nextVesre = request.RequestedVerse + one
	} else {
		nextVesre = one
	}

	responce := models.VerseResponce{
		NextVerse: nextVesre,
		Text:      sliceVerses[request.RequestedVerse-one],
	}
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "GetSong"}).Infof("success. nextVesre:%v", nextVesre)

	return responce, nil
}

func (s *ApplicationServices) GetAllSongs(ctx context.Context, requestSongFilter models.RequestSongsFilter) (*models.ResponceSongs, error) {
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "GetAllSongs"}).Infof("request:%+v", requestSongFilter)

	songsFilterDBModel := models.NewSongsFilterDBModel(requestSongFilter)

	sliceSongs, err := s.repository.GetAllSongs(ctx, songsFilterDBModel)
	if err != nil || len(sliceSongs) == 0 || sliceSongs == nil {
		s.log.WithError(err).Error("failed to get all songs from repository")
		return nil, err
	}
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "GetAllSongs"}).Infof("success. len sliceSongs:%+v", len(sliceSongs))

	return models.NewResponceSongs(sliceSongs), nil
}

func (s *ApplicationServices) DeleteSong(ctx context.Context, id int) error {
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "DeleteSong"}).Infof("entry: id %v", id)

	err := s.repository.DeleteSongFromRepository(ctx, id)
	if err != nil {
		s.log.WithError(err).Error("failed to delete song from repository")
		return err
	}
	s.log.WithFields(logrus.Fields{"layer": "services", "op": "DeleteSong"}).Infof("success: id %v", id)

	return nil
}
