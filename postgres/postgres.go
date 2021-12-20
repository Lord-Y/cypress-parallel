// Package postgres stand to manage db init and ping
package postgres

import (
	"database/sql"
	"embed"
	"net/http"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/rs/zerolog/log"
)

// InitDB permit to initialiaze or migrate databases
func InitDB(source embed.FS) {
	src, err := httpfs.New(http.FS(source), "sql/postgres")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to get embedded sql directory")
		return
	}
	log.Debug().Msg("starting db initialization/migration")
	log.Debug().Msgf("Use db sql driver %s", "postgres")
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to connect to DB")
		return
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msgf("could not ping DB: %s", err.Error())
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not start sql migration with error msg: %s", err.Error())
		return
	}
	m, err := migrate.NewWithInstance(
		"httpfs",
		src,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Migration failed: %s", err.Error())
		return
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msgf("An error occurred while syncing the database with error msg: %s", err.Error())
		return
	}
	log.Info().Msg("Database migrated successfully")
}

// Ping permit to ping sql instance
func Ping() (b bool) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while connecting to DB")
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while pinging DB")
		return
	}
	return true
}
