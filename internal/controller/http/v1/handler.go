package v1

import (
	"context"

	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	groupName = "group"
	songName  = "song"
)

type SongService interface {
	AddSong(context.Context, models.SongRequest, []byte) (int, error)
	DeleteSong(context.Context, int) error
	UpdateSong(context.Context, models.SongUpdate) error
}

type Handler struct {
	log        *logrus.Logger
	configHTTP config.HTTP
	service    SongService
}

func NewHandler(log *logrus.Logger, service SongService, cfg config.HTTP) *Handler {
	return &Handler{
		log:        log,
		configHTTP: cfg,
		service:    service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	songs := router.Group("/songs")
	{
		songs.POST("/add", h.addSong)
		songs.PUT("/update", h.updateSong)
		// auth.GET("/get/:id", h.getClient)
		songs.DELETE("/delete/:id", h.deleteSong)
		// auth.GET("/statistic", h.getStatistic)
	}

	return router
}
