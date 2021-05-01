// Package postgres stand to manage db init and ping
package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/rs/zerolog/log"
)

// InitDB permit to initialiaze or migrate databases
func InitDB() {
	log.Debug().Msg("starting db initialization/migration")
	fileDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Not able to get current directory")
	}
	log.Debug().Msgf("Use db sql driver %s", "postgres")
	sqlDIR := fmt.Sprintf("file://%s%s", fileDir, "/sql/postgres")
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
		log.Fatal().Err(err).Msgf("could not ping DB: %v", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not start sql migration with error msg: %v", err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		sqlDIR,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Migration failed: %v", err)
		return
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msgf("An error occurred while syncing the database with error msg: %v", err)
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
	defer db.Close()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while connecting to DB")
		return
	}
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while pinging DB")
		return
	}
	return true
}
