package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Lord-Y/cypress-parallel-api/hooks"
	customLogger "github.com/Lord-Y/cypress-parallel-api/logger"
	"github.com/Lord-Y/cypress-parallel-api/postgres"
	"github.com/Lord-Y/cypress-parallel-api/routers"
	"github.com/rs/zerolog/log"
)

// init func
func init() {
	customLogger.SetLoggerLogLevel()

	if strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_DB_URI")) == "" {
		msg := "CYPRESS_PARALLEL_DB_URI environment variable must be set"
		log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
		return
	}

	postgres.InitDB()
}

func main() {
	var srv *http.Server
	router := routers.SetupRouter()

	appPort := strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_PORT"))
	if appPort != "" {
		srv = &http.Server{
			Addr:    fmt.Sprintf(":%s", appPort),
			Handler: router,
		}
		log.Info().Msgf("Starting server on port %s", appPort)
	} else {
		srv = &http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		log.Info().Msg("Starting server on port 8080")
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Startup failed")
		}
	}()

	go queued()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server shutted down abruptly")
	}
	log.Info().Msg("Server exited successfully")
}

func queued() {
	for range time.Tick(30 * time.Second) {
		hooks.Queued()
	}
}
