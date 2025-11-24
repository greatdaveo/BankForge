package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(rateHandler http.HandlerFunc) http.Handler {
	r := chi.NewRouter()

	// To test the root route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" XChange is working!!!"))
	})

	r.Get("/convert", rateHandler)

	return r
}
