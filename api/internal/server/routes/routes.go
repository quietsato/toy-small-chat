package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/quietsato/toy-small-chat/api/internal/di"
)

func Setup(r *chi.Mux, dic *di.Container) {
	tokenAuth := dic.Auth.Middleware.GetTokenAuthForMiddleware()

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Post("/login", login(dic))
		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", createAccount(dic))
		})
	})
	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(accountCtx)
		// Room
		r.Route("/rooms", func(r chi.Router) {
			r.Get("/", getRooms(dic))
			r.Post("/", createRoom(dic))
		})
		// Message
		r.Route("/rooms/{roomID}/messages", func(r chi.Router) {
			r.Use(roomCtx)
			r.Get("/", getMessages(dic))
			r.Post("/", createMessage(dic))
		})
	})
}
