# goauth

Work-in-progress middleware module that can be used with Go-Chi as well as other frameworks.

Offers token and cookie-based authentication

It should be straight-forward to use on specific routes:

```
import (
        "net/http"

        "github.com/amonaco/goauth/lib/auth"
        "github.com/go-chi/chi"
)

func main() {

        r := chi.NewRouter()
        r.Use(middleware.Logger)

	resource := "/healthcheck" 
        r.Get(resource, func(w http.ResponseWriter, r *http.Request) {

                // Auth middleware here
                r.Use(auth.Middleware(resource))
                w.Write([]byte("OK"))
        })
}

```
