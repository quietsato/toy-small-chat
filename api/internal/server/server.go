package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/quietsato/toy-small-chat/api/internal/di"
	"github.com/quietsato/toy-small-chat/api/internal/server/routes"

	slogchi "github.com/samber/slog-chi"
)

func New(dic *di.Container) *chi.Mux {
	r := chi.NewRouter()

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:    func(r *http.Request, origin string) bool { return true }, // Not for production, allow all origins
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{},
		AllowCredentials:   false,
		MaxAge:             86400, // 24h
		OptionsPassthrough: false,
		Debug:              false,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(slogchi.New(slog.Default()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	routes.Setup(r, dic)

	return r
}
