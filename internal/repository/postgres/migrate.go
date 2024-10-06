package postgres

import (
	"log"

	"github.com/andrew-nino/em_songs/config"
	go_migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Returns a new Migrate instance
func NewMigration(c *config.PG) *go_migrate.Migrate {
	m, err := go_migrate.New(
		"file://schema",
		"postgres://"+c.Username+":"+c.Password+"@"+c.Host+":"+c.Port+"/"+c.DBName+"?sslmode="+c.SSLMode+"")

	if err != nil {
		log.Fatalf("failed to migrate db: %s", err.Error())
	}

	return m
}
