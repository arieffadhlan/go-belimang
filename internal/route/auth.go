package route

import (
	"belimang/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, h handlers.AuthHandler) {
	r.Post("/admin/login", func(w http.ResponseWriter, r *http.Request) { h.SignIn(w, r, "admin") })
	r.Post("/users/login", func(w http.ResponseWriter, r *http.Request) { h.SignIn(w, r, "users") })
	r.Post("/admin/register", func(w http.ResponseWriter, r *http.Request) { h.SignUp(w, r, "admin") })
	r.Post("/users/register", func(w http.ResponseWriter, r *http.Request) { h.SignUp(w, r, "users") })
}
