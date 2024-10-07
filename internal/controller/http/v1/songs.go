package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/go-playground/validator"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addSong(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("Error reading request body")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
	}

	requestModel := models.SongRequest{}
	err = bindAndValidateRequest(body, &requestModel)
	if err != nil {
		h.log.Error("failed to bind JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	ctx := c.Request.Context()
	rawAnswer, err := getSongFromExternal(ctx, h.configHTTP, body)
	if err != nil {
		h.log.Error("failed to get song from external service: ", err)
		c.JSON(http.StatusFailedDependency, gin.H{"message": "failed to get song from external service"})
		return
	}

	id, err := h.service.AddSong(ctx, requestModel, rawAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to add song to db"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func bindAndValidateRequest(body []byte, request *models.SongRequest) error {

	if err := json.Unmarshal(body, request); err != nil {
		return err
	}

	var validate = validator.New()

	if err := validate.Struct(request); err != nil {
		return err
	}

	return nil
}

func getSongFromExternal(ctx context.Context, cfg config.HTTP, requestBody []byte) ([]byte, error) {

	// TODO: добавить retry
	tr := &http.Transport{
		MaxIdleConnsPerHost:   cfg.MaxIdleConns,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.ExternalURL, bytes.NewReader(requestBody))
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
