package v1

import (
	"context"

	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SongService interface {
	AddSong(context.Context, models.SongRequest, []byte) (int, error)
}

type Handler struct {
	log      *logrus.Logger
	confHTTP config.HTTP
	service  SongService
}

func NewHandler(log *logrus.Logger, service SongService, cfg config.HTTP) *Handler {
	return &Handler{
		log:      log,
		confHTTP: cfg,
		service:  service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	songs := router.Group("/songs")
	{
		songs.POST("/add", h.addSong)
		// auth.PUT("/update", h.updateClient)
		// auth.GET("/get/:id", h.getClient)
		// auth.DELETE("/delete/:id", h.deleteClient)
		// auth.GET("/statistic", h.getStatistic)
	}

	return router
}
