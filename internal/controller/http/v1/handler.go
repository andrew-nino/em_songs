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
	UpdateSong(context.Context, models.SongUpdate) error
	GetSong(context.Context, models.VerseRequest) (models.VerseResponce, error)
	GetAllSongs(context.Context, models.RequestSongsFilter) ([]models.ResponceSongs, error)
	DeleteSong(context.Context, int) error
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
		songs.GET("/get_one", h.getSong)
		songs.GET("/get_all", h.getAllSongs)
		songs.DELETE("/delete/:id", h.deleteSong)
	}

	return router
}
