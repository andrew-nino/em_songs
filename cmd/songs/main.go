package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrew-nino/em_songs/config"
	"github.com/andrew-nino/em_songs/internal/app"
)

func main() {

	cfg := config.NewConfig()

	log := SetLogrus(cfg.Log.Level)

	application := app.NewApplication(log, cfg.HTTP.Port, cfg)

	go application.HTTPServer.MustRun()

	log.Print("App " + cfg.App.Name + " version: " + cfg.App.Version + " started")

	// Waiting signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	application.HTTPServer.Shutdown(context.Background())

	log.Print("App " + cfg.App.Name +  " stopped")
}
