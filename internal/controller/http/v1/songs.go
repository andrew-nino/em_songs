package v1

import (
	"context"
	"io"
	"net/http"

	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/go-playground/validator"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addSong(c *gin.Context) {
	ctx := c.Request.Context()
	request := models.SongRequest{}

	err := bindAndValidate(c, &request)
	if err != nil {
		h.log.Error("failed to bind JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	rawAnswer, err := getSong(ctx, c, h.confHTTP)
	if err != nil {
		h.log.Error("failed to get song from external service: ", err)
		c.JSON(http.StatusFailedDependency, gin.H{"message": "failed to get song from external service"})
		return
	}

	id, err := h.service.AddSong(ctx, request, rawAnswer)
	if err != nil {
		h.log.Error("failed to add song to db: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to add song to db"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func bindAndValidate(c *gin.Context, request *models.SongRequest) error {

	if err := c.BindJSON(request); err != nil {
		return err
	}

	var validate = validator.New()

	if err := validate.Struct(request); err != nil {
		return err
	}

	return nil
}

func getSong(ctx context.Context, c *gin.Context, cfg config.HTTP) ([]byte, error) {

	body := c.Request.Body
	defer body.Close()

	rawAnswer, err := getSongFromExternal(ctx, cfg, body)

	return rawAnswer, err
}

func getSongFromExternal(ctx context.Context, cfg config.HTTP, requestBody io.ReadCloser) ([]byte, error) {

	// TODO: добавить retry
	tr := &http.Transport{
		MaxIdleConnsPerHost:   cfg.MaxIdleConns,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.ExternalURL, requestBody)
	if err != nil {
		return nil, err
	}

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
