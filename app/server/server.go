package server

import (
	"fmt"
	"net/http"
	"time"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func NewServer(address string, port int) *http.Server {
	return &http.Server{
        Addr: fmt.Sprintf("%s:%d", address, port),
        Handler: routes(),
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 60 * time.Second,
        IdleTimeout: 30 * time.Second,
	}
}

func routes() chi.Router {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.Route("/internal", func(r chi.Router) {
		r.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
            writer.Header().Set("Content-Type", "text/plain")
            writer.WriteHeader(http.StatusOK)
            _, _ = writer.Write([]byte("pong"))
		})

	})

    return r
}