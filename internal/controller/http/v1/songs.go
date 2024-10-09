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

	"github.com/gin-gonic/gin"
)

func (h *Handler) addSong(c *gin.Context) {

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
	rawAnswer, err := getSongFromExternal(ctx, h.configHTTP, requestModel)
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

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) updateSong(c *gin.Context) {

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

	c.JSON(http.StatusOK, gin.H{"successful update song id": updateSongModel.SongID})
}

func (h *Handler) getSong(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Error reading request body")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
	}

	getVerseModel := models.VerseRequest{}
	err = bindAndValidateRequest(body, &getVerseModel)
	if err != nil {
		h.log.Error("failed to bind JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	verse, err := h.service.GetSong(ctx, getVerseModel)
	if err != nil {
		if errors.Is(err, ErrNoSongFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to getting verse"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to getting verse"})
		return
	}

	c.JSON(http.StatusOK, verse)
}

func (h *Handler) deleteSong(c *gin.Context) {

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

	c.JSON(http.StatusOK, gin.H{"successful deletion of id ": id})
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
