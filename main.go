// Sample healthcheck endpoint service using the middleware
package main

import (
	"log"
	"net/http"

	"github.com/amonaco/goauth/lib/config"
	"github.com/amonaco/goauth/lib/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	config.Read("config/config.yml")
	conf := config.Get()

	router := chi.NewRouter()

	// Redis
	// cache.Start()

	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		router.Use(auth.Middleware)
		w.Write([]byte("OK"))
	})

	log.Printf("Environment: %v", conf.Environment)
	log.Printf("Listening on %v", conf.Listen)
	log.Fatal(http.ListenAndServe(conf.Listen, router))

}
