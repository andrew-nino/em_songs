package v1

import (
	"context"
	"strings"
	"time"

	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	groupName  = "group"
	songName   = "song"
	defaultURL = "http://localhost:8080"
)

type SongService interface {
	AddSong(context.Context, models.SongRequest, []byte) (int, error)
	UpdateSong(context.Context, models.SongUpdate) error
	GetSong(context.Context, models.VerseRequest) (models.VerseResponce, error)
	GetAllSongs(context.Context, models.RequestSongsFilter) (*models.ResponceSongs, error)
	DeleteSong(context.Context, int) error
}

type Handler struct {
	log     *logrus.Logger
	config  config.Config
	service SongService
}

func NewHandler(log *logrus.Logger, service SongService, cfg *config.Config) *Handler {
	return &Handler{
		log:     log,
		config:  *cfg,
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	if h.config.Gin.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	var allowOrigins = make([]string, 0)

	if !h.config.Gin.AllowAllOrigins {
		h.log.Infof("CORS enabled")

		allowOrigins = strings.Split(h.config.Gin.AllowUrls, ",")
		allowOrigins = append(allowOrigins, defaultURL)
		h.log.Infof("allowOriginsURLs: %+v", allowOrigins)
	}

	corsConfig := cors.Config{
		AllowAllOrigins:  h.config.Gin.AllowAllOrigins,
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	songs := router.Group("/songs")
	v1 := songs.Group("/v1")
	{
		v1.POST("/add", h.addSong)
		v1.PUT("/update", h.updateSong)
		v1.GET("/get_one", h.getSong)
		v1.GET("/get_all", h.getAllSongs)
		v1.DELETE("/delete/:id", h.deleteSong)
	}

	return router
}
