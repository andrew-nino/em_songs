package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/andrew-nino/em_songs/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	groups_table     = "groups"
	songs_table      = "songs"
	group_song_table = "group_song"
)

type Postgres struct {
	log          *logrus.Logger
	connAttempts int
	connTimeout  time.Duration
	db           *sqlx.DB
}

func New(log *logrus.Logger, cfg config.PG) *Postgres {

	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	// Migrates running
	log.Info("Migrates running...")
	m := NewMigration(cfg)
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Print("no migrations to apply")
		} else {
			log.Fatalf("failed to apply migrations: %s", err.Error())
		}
	}

	return &Postgres{
		log:          log,
		connAttempts: cfg.ConnAttempts,
		connTimeout:  cfg.ConnTimeout,
		db:           db,
	}
}

// Causes the database to open and checks the connection. If the connection is established, returns a pointer to the database.
// Returns an error if the database has not opened or there is no connection.
func NewPostgresDB(cfg config.PG) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	for cfg.ConnAttempts > 0 {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Postgres is trying to connect, attempts left: %d\n", cfg.ConnAttempts)
		time.Sleep(cfg.ConnTimeout)
		cfg.ConnAttempts--
	}

	return db, err
}
