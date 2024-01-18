package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF_Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// a heartbeat route, to ensure things are up
	mux.Use(middleware.Heartbeat("/ping"))

	// this route is just to ensure things work, and is never
	// used after that
	mux.Get("/", app.Broker)

	mux.Post("/", app.Broker)

	// a route for everything
	mux.Post("/handle", app.HandleSubmission)

	// grpc route
	mux.Post("/log-grpc", app.LogViaGRPC)

	return mux
}
