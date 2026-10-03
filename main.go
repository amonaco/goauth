package main

import (
	"log"
	"net/http"

	"github.com/amonaco/goauth/lib/auth"
	"github.com/amonaco/goauth/lib/config"
	"github.com/amonaco/goauth/lib/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.Read("config/config.yml")
	conf := config.Get()

	router := chi.NewRouter()
	router.Use(middleware.Middleware)

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.Get("/secure", func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserID(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		w.Write([]byte("authenticated user: " + string(rune(userID))))
	})

	log.Printf("Environment: %v", conf.Environment)
	log.Printf("Listening on %v", conf.Listen)
	log.Fatal(http.ListenAndServe(conf.Listen, router))
}
