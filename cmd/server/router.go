package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//NewRouter configure router with api endpoints
func NewRouter(c *Config, s *Store) *chi.Mux {
	r := chi.NewRouter()
	r.Use(JSONCtx)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.MethodNotAllowed(NotAllowedHandler())
	r.NotFound(NotFoundHandler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", PingHandler())

		r.Post("/login", LoginHandler(c, s))

		r.Route("/keys", func(r chi.Router) {
			if c.Authorization {
				r.Use(Authorization(s))
			}

			r.Get("/{key}/values", GetHandler(s))
			r.Get("/{key}/values/{index}", GetIndexHandler(s))
			r.Post("/", SetHandler(s))

			r.Get("/{pattern}", KeysHandler(s))
			r.Delete("/{key}", RemoveHandler(s))

			r.Get("/{key}/expires", GetExpiresHandler(s))
			r.Post("/{key}/expires", SetExpiresHandler(s))
		})
	})
	return r
}
