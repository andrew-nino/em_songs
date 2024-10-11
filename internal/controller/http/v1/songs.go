package v1

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addSong(c *gin.Context) {
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "addSong"}).Infof("entry: IP %+v", c.RemoteIP())

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Error reading request body")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
	}

	requestModel := models.SongRequest{}
	err = bindAndValidateRequest(body, &requestModel)
	if err != nil {
		h.log.Error("failed to bind JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	rawAnswer, err := getSongFromExternal(ctx, h.config.HTTP, requestModel)
	if err != nil {
		h.log.Error("failed to get song from external service: ", err)
		c.JSON(http.StatusFailedDependency, gin.H{"message": "failed to get song"})
		return
	}

	id, err := h.service.AddSong(ctx, requestModel, rawAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to add song"})
		return
	}
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "addSong"}).Infof("success: IP %+v", c.RemoteIP())

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) updateSong(c *gin.Context) {
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "updateSong"}).Infof("entry: IP %+v", c.RemoteIP())

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Error reading request body")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
	}

	updateSongModel := models.SongUpdate{}
	err = bindAndValidateRequest(body, &updateSongModel)
	if err != nil {
		h.log.Error("failed to bind JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdateSong(ctx, updateSongModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update song"})
		return
	}
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "updateSong"}).Infof("success: IP %+v", c.RemoteIP())

	c.JSON(http.StatusOK, gin.H{"successful update song id": updateSongModel.SongID})
}

func (h *Handler) getSong(c *gin.Context) {
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "getSong"}).Infof("entry: IP %+v", c.RemoteIP())

	id := c.Query("id")
	requestedVerse := c.Query("requestedVerse")

	getVerseModel, err := validatingGetSong(id, requestedVerse)
	if err != nil {
		h.log.Error("failed to bind JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	verse, err := h.service.GetSong(ctx, *getVerseModel)
	if err != nil {
		if errors.Is(err, ErrNoSongFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to getting verse"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to getting verse"})
		return
	}
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "getSong"}).Infof("success: IP %+v", c.RemoteIP())

	c.JSON(http.StatusOK, verse)
}

func (h *Handler) getAllSongs(c *gin.Context) {
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "getAllSongs"}).Infof("entry: IP %+v", c.RemoteIP())

	limit := c.Query("limit")
	offset := c.Query("offset")
	group := c.Query("group")
	song := c.Query("song")

	requestSongFilter, err := validatingQuerysFilter(limit, offset, group, song)
	if err != nil {
		h.log.Error("failed to validate params: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	sliceResponceSongs, err := h.service.GetAllSongs(ctx, *requestSongFilter)
	if err != nil || sliceResponceSongs == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get songs"})
		return
	}
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "getAllSongs"}).Infof("success: IP %+v", c.RemoteIP())

	c.JSON(http.StatusOK, *sliceResponceSongs)
}

func (h *Handler) deleteSong(c *gin.Context) {
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "deleteSong"}).Infof("entry: IP %+v", c.RemoteIP())

	ctx := c.Request.Context()
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error("invalid id: ", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	err = h.service.DeleteSong(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete song"})
		return
	}
	h.log.WithFields(logrus.Fields{"layer": "handler", "op": "deleteSong"}).Infof("success: IP %+v", c.RemoteIP())

	c.JSON(http.StatusOK, gin.H{"successful deletion of id ": id})
}

func validatingGetSong(id, requestedVerse string) (*models.VerseRequest, error) {

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	vesreInt, err := strconv.Atoi(requestedVerse)
	if err != nil {
		return nil, err
	}

	getVerseModel := &models.VerseRequest{
		ID:             int64(idInt),
		RequestedVerse: int64(vesreInt),
	}

	var validate = validator.New()
	if err := validate.Struct(getVerseModel); err != nil {
		return nil, err
	}
	return getVerseModel, nil
}

func validatingQuerysFilter(params ...string) (*models.RequestSongsFilter, error) {

	limit, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}
	offset, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, err
	}

	requestSongFilter := models.NewRequestSongFilter(int64(limit), int64(offset), params[2], params[3])

	var validate = validator.New()
	if err := validate.Struct(requestSongFilter); err != nil {
		return nil, err
	}
	return &requestSongFilter, nil
}

func bindAndValidateRequest(body []byte, model interface{}) error {

	if err := json.Unmarshal(body, model); err != nil {
		return err
	}

	var validate = validator.New()

	if err := validate.Struct(model); err != nil {
		return err
	}

	return nil
}

func getSongFromExternal(ctx context.Context, cfg config.HTTP, requestModel models.SongRequest) ([]byte, error) {

	// TODO: добавить retry
	tr := &http.Transport{
		MaxIdleConnsPerHost:   cfg.MaxIdleConns,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.ExternalURL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set(groupName, requestModel.Group)
	query.Set(songName, requestModel.Song)
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
