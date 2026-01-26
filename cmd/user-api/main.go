package main

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jlargs64/gameframe-microservices/internal/server"
)

func main() {
	// Set up handlers
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Info("starting the user service...")
	// Create server
	server, err := server.New(
		server.WithPort("8080"),
		server.WithEnv("dev"),
		server.WithHandler(r),
	)
	if err != nil {
		log.Fatal("could not  the server")
	}
	log.Info("user service started!")
	// Start the server
	server.Start()
	log.Info("the user service was stopped")
}
