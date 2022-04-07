// Sample healthcheck endpoint service using the middleware
package main

import (
	"net/http"

	"github.com/amonaco/goauth/lib/auth"
	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {

		// Auth middleware here
		r.Use(auth.Middleware)
		// TODO: See of passing arguments to the middleware here (like access object)
		w.Write([]byte("OK"))
	})
}
