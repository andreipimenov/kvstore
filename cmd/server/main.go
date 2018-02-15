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

	s := NewStore(store.New())

	r := chi.NewRouter()
	r.Use(JSONCtx)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.MethodNotAllowed(NotAllowedHandler())
	r.NotFound(NotFoundHandler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", PingHandler())
		r.Post("/values", SetHandler(s))
		r.Get("/values/{key}", GetHandler(s))
		r.Get("/values/{key}/{index}", GetIndexHandler(s))
		r.Delete("/keys/{key}", RemoveHandler(s))
		r.Get("/keys/{pattern}", KeysHandler(s))
	})

	log.Printf("Start listening on port %d", c.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
	if err != nil {
		log.Fatal(err)
	}
}
