package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andreipimenov/kvstore/config"
	"github.com/andreipimenov/kvstore/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	c, err := NewConfig(config.New("etc/config.json"))
	if err != nil {
		log.Fatal(err)
		return
	}

	s := NewStore(store.New(c.DumpFile, c.DumpInterval))

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

	log.Printf("Start listening on port %d", c.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
	if err != nil {
		log.Fatal(err)
	}
}
