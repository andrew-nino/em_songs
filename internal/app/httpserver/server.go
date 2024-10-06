package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/andrew-nino/em_songs/config"
	v1 "github.com/andrew-nino/em_songs/internal/controller/http/v1"
	"github.com/andrew-nino/em_songs/internal/service"
	"github.com/sirupsen/logrus"
)

type Server struct {
	log        *logrus.Logger
	handler    *v1.Handler
	httpServer *http.Server
	port       string
}

func New(log *logrus.Logger, port string, services *service.ApplicationServices, cfg config.HTTP) *Server {

	handler := v1.NewHandler(log, services, cfg)

	httpServer := &http.Server{
		Addr:           ":" + port,
		Handler:        handler.InitRoutes(),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    12 * time.Second,
		WriteTimeout:   12 * time.Second,
	}

	return &Server{
		log:        log,
		handler:    handler,
		httpServer: httpServer,
		port:       port,
	}
}

// Configure with the necessary parameters and start the server.
func (s *Server) MustRun() {

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.log.Fatalf("HTTP server failed to start: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.log.Infof("HTTP server shutdown with error: %v", err)
		return err
	}
	return nil
}
