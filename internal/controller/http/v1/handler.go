package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GroupService interface {
}

type SongsService interface {
}

type Handler struct {
	log    *logrus.Logger
	groups GroupService
	songs  SongsService
}

func NewHandler(log *logrus.Logger, groups GroupService, songs SongsService) *Handler {
	return &Handler{
		log:    log,
		groups: groups,
		songs:  songs,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	_ = router.Group("/info")
	{
		// auth.POST("/add", h.addClient)
		// auth.PUT("/update", h.updateClient)
		// auth.GET("/get/:id", h.getClient)
		// auth.DELETE("/delete/:id", h.deleteClient)
		// auth.GET("/statistic", h.getStatistic)
	}

	return router
}
