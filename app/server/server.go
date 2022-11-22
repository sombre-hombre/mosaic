package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sombre-hombre/mosaic/app/mosaic"
	"github.com/sombre-hombre/mosaic/app/server/models"
	"github.com/sombre-hombre/mosaic/app/tiles"
	"image"
	"image/jpeg"
	"net/http"
	"time"
)

// NewServer creates new API server
func NewServer(address string, port int) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", address, port),
		Handler:      routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

func routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Route("/internal", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("pong"))
		})
		r.Mount("/pprof", middleware.Profiler())
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.URLFormat)
		r.Route("/api/v1", func(apiRouter chi.Router) {
			apiRouter.Use(middleware.AllowContentType("image/png", "image/jpeg"))
			apiRouter.Post("/libraries/{ID}/mosaics", mosaicsHandler)
		})
	})

	return router
}

// POST /api/v1/libraries/{ID}/mosaics
func mosaicsHandler(w http.ResponseWriter, r *http.Request) {
	libID := chi.URLParam(r, "ID")
	if libID == "" {
		render.Render(w, r, models.ErrNotFound)
		return
	}

	// TODO: cache libraries
	lib, err := tiles.NewLibrary(fmt.Sprintf("../../img/tiles/%s", libID), 50, tiles.ColorDistanceRedmean)
	if err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}

	sourceImg, _, err := image.Decode(r.Body)
	if err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}

	m, err := mosaic.Create(sourceImg, *lib)
	if err != nil {
		render.Render(w, r, models.ErrServer(err))
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	_ = jpeg.Encode(w, m, nil) // TODO: error handling
}
