// Package server allows for the
// easy construction of HTTP servers
package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
)

type Server struct {
	Port    string
	Env     string
	Handler http.Handler
}

const (
	defaultPort = "8080"
	devEnv      = "dev"
	prodEnv     = "prod"
	defaultEnv  = devEnv
)

// New creates a new server with the functional options pattern
func New(options ...func(*Server) error) (*Server, error) {
	server := &Server{
		defaultPort,
		defaultEnv,
		nil,
	}
	// Apply functional options
	for _, opt := range options {
		if err := opt(server); err != nil {
			return nil, err
		}
	}
	return server, nil
}

// Start the HTTP server
func (s *Server) Start() {
	httpServer := http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Handler,
	}

	// Setup shutdown context to handle graceful shutdowns
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Start the server in a separate goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("could not start the server", "err", err)
			}
		}
	}()

	// Listen for shutdown signals
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal("the server could not shutdown gracefully", "err", err)
	}
}

// WithPort specifies the port for a new API server
// and returns an error when the param is invalid
func WithPort(port string) func(*Server) error {
	return func(s *Server) error {
		if len(port) == 0 {
			return errors.New("the port must be specified and not empty")
		}
		s.Port = port
		return nil
	}
}

// WithEnv specifies the env for a new API server
// and returns an error when the param is invalid
func WithEnv(env string) func(*Server) error {
	return func(s *Server) error {
		// Ensure the env is a valid option
		if env != devEnv && env != prodEnv {
			return errors.New("the server env must be dev or env")
		}
		s.Env = env
		return nil
	}
}

// WithHandler specifies the API handlers for a the API server
// and returns an error when the param is invalid
func WithHandler(handler http.Handler) func(*Server) error {
	return func(s *Server) error {
		s.Handler = handler
		return nil
	}
}
