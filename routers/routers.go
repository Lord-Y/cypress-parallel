// Package routers expose all routes of the api
package routers

import (
	"embed"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Lord-Y/cypress-parallel/annotations"
	"github.com/Lord-Y/cypress-parallel/environments"
	"github.com/Lord-Y/cypress-parallel/executions"
	"github.com/Lord-Y/cypress-parallel/health"
	"github.com/Lord-Y/cypress-parallel/hooks"
	customLogger "github.com/Lord-Y/cypress-parallel/logger"
	"github.com/Lord-Y/cypress-parallel/projects"
	"github.com/Lord-Y/cypress-parallel/teams"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	customLogger.SetLoggerLogLevel()
}

//go:generate rm -rf ui/* && cp -r ../ui/dist ui
//go:embed ui/dist/assets
var assets embed.FS

//go:embed ui
var ui_ embed.FS

// SetupRouter func handle all routes of the api
func SetupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	requestID := tools.RandStringInt(32)

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	router := gin.New()
	router.Use(gin.Recovery())
	router.RedirectTrailingSlash = true

	router.Use(
		logger.SetLogger(
			logger.WithUTC(true),
			logger.WithLogger(
				func(c *gin.Context, w io.Writer, d time.Duration) zerolog.Logger {
					return zerolog.New(os.Stdout).
						With().
						Timestamp().
						Str("requestId", requestID).
						Int("status", c.Writer.Status()).
						Str("method", c.Request.Method).
						Str("path", c.Request.URL.Path).
						Str("ip", c.ClientIP()).
						Dur("latency", d).
						Str("user_agent", c.Request.UserAgent()).
						Logger()
				},
			),
		),
	)
	headerHandler := func(c *gin.Context) {
		if c.GetHeader("X-Request-Id") == "" {
			c.Request.Header.Set("X-Request-Id", requestID)
			c.Next()
		}
	}
	router.Use(headerHandler)

	// disable during unit testing
	if os.Getenv("CYPRESS_PARALLEL_PROMETHEUS") != "" {
		p := ginprometheus.NewPrometheus("http")
		p.SetListenAddress(":9101")
		p.Use(router)
	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", health.Health)

		v1.POST("/teams", teams.Create)
		v1.GET("/teams/:teamId", teams.Read)
		v1.GET("/teams/list", teams.List)
		v1.GET("/teams/all", teams.All)
		v1.PUT("/teams", teams.Update)
		v1.DELETE("/teams/:teamId", teams.Delete)
		v1.GET("/teams/search", teams.Search)

		v1.POST("/projects", projects.Create)
		v1.GET("/projects/:projectId", projects.Read)
		v1.GET("/projects/list", projects.List)
		v1.GET("/projects/all", projects.All)
		v1.PUT("/projects", projects.Update)
		v1.DELETE("/projects/:projectId", projects.Delete)
		v1.GET("/projects/search", projects.Search)

		v1.POST("/environments", environments.Create)
		v1.PUT("/environments", environments.Update)
		v1.GET("/environments/list", environments.List)
		v1.GET("/environments/list/by/projectid/:projectId", environments.ListByProjectID)
		v1.GET("/environments/:environmentId", environments.Read)
		v1.DELETE("/environments/:environmentId", environments.Delete)
		v1.GET("/environments/search", environments.Search)

		v1.POST("/annotations", annotations.Create)
		v1.PUT("/annotations", annotations.Update)
		v1.GET("/annotations/list", annotations.List)
		v1.GET("/annotations/list/by/projectid/:projectId", annotations.ListByProjectID)
		v1.GET("/annotations/:annotationId", annotations.Read)
		v1.DELETE("/annotations/:annotationId", annotations.Delete)
		v1.GET("/annotations/search", annotations.Search)

		v1.POST("/hooks/launch/plain", hooks.Plain)

		v1.GET("/executions/list", executions.List)
		v1.GET("/executions/list/by/uniqid/:uniqId", executions.UniqID)
		v1.POST("/executions/update", executions.UpdateResultExecution)
		v1.GET("/executions/:executionId", executions.Read)
		v1.GET("/executions/search", executions.Search)
	}

	router.Static("/ui/assets", "ui/dist/assets")

	// done like that to avoid trailing slash
	router.GET("/ui/logo.png", func(c *gin.Context) {
		f, _ := ui_.ReadFile("ui/dist/logo.png")
		c.Data(
			http.StatusOK,
			"image/png",
			f,
		)
	})

	// done like that to avoid trailing slash
	router.GET("/ui/favicon.ico", func(c *gin.Context) {
		f, _ := ui_.ReadFile("ui/dist/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			f,
		)
	})

	ui := router.Group("/ui", func(c *gin.Context) {
		f, _ := ui_.ReadFile("ui/dist/index.html")
		c.Data(
			http.StatusOK,
			"text/html",
			f,
		)
	})
	if ui.BasePath() == "/ui" {
		router.NoRoute(func(c *gin.Context) {
			f, _ := ui_.ReadFile("ui/dist/index.html")
			c.Data(
				http.StatusOK,
				"text/html",
				f,
			)
		})
	}

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/ui/")
	})

	return router
}
