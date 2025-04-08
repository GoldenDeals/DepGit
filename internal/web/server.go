// Package web provides the web server and API handlers for the DepGit application.
package web

import (
	"context"
	"net/http"
	"time"

	"github.com/GoldenDeals/DepGit/internal/gen/api"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var serverLogger = logger.New("web-server")

// Config holds the configuration for the web server
type Config struct {
	Address     string
	StaticDir   string
	APIBasePath string
}

// Server represents the web server for the DepGit application
type Server struct {
	config  Config
	echo    *echo.Echo
	handler api.ServerInterface
}

// New creates a new web server with the given configuration and API handler
func New(config Config, handler api.ServerInterface) *Server {
	e := echo.New()

	// Configure middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Create server
	s := &Server{
		config:  config,
		echo:    e,
		handler: handler,
	}

	// Set up routes
	s.setupRoutes()

	return s
}

// setupRoutes configures the routes for the web server
func (s *Server) setupRoutes() {
	// Serve static files
	s.echo.Static("/", s.config.StaticDir)

	// Set up API routes
	api.RegisterHandlersWithBaseURL(s.echo, s.handler, s.config.APIBasePath)

	// SPA fallback - serve index.html for any unmatched routes
	s.echo.GET("*", func(c echo.Context) error {
		return c.File(s.config.StaticDir + "/index.html")
	})
}

// Start starts the web server and listens for incoming connections
func (s *Server) Start(ctx context.Context) error {
	// Start server in a goroutine
	go func() {
		if err := s.echo.Start(s.config.Address); err != nil && err != http.ErrServerClosed {
			serverLogger.WithError(err).Fatal("Failed to start web server")
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.echo.Shutdown(shutdownCtx)
}

// Close shuts down the web server
func (s *Server) Close() error {
	return s.echo.Close()
}
