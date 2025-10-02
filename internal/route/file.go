package route

import (
	"belimang/internal/handlers"
	"belimang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterFileRoutes(r chi.Router, h *handlers.FileHandler) {
	r.Group(func(g chi.Router) {
		g.Use(middleware.Protected(true))
		g.Post("/image", h.UploadFile)
	})
}
