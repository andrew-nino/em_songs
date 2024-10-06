package app

import (
	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/app/httpserver"
	"github.com/andrew-nino/em_songs/internal/repository/postgres"
	"github.com/andrew-nino/em_songs/internal/service"
	"github.com/sirupsen/logrus"
)

type App struct {
	HTTPServer *httpserver.Server
}

func NewApplication(log *logrus.Logger, port string, cfg *config.Config) *App {

	repository := postgres.New(log, cfg.PG)
	services := service.New(log, repository)
	server := httpserver.New(log, port, services, cfg.HTTP)

	return &App{
		HTTPServer: server,
	}
}
