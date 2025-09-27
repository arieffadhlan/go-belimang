package route

import (
	"belimang/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterFileRoutes(r chi.Router, h *handlers.FileHandler) {
	r.Group(func(g chi.Router) {
		// g.Use(middleware.Protected)
		// g.Post("/image", h.CreateProduct)
	})
}
