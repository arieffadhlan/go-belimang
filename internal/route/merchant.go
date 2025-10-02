package route

import (
	"belimang/internal/handlers"
	"belimang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterMerchantRoutes(r chi.Router, h *handlers.MerchantHandler) {
	r.Group(func(g chi.Router) {
		g.Use(middleware.Protected(true))

		g.Get("/admin/merchants", h.GetAllMerchant)
		// g.Get("/admin/merchants/:merchantId/items", h.GetAllMerchantItems)

		g.Post("/admin/merchants", h.CreateMerchant)
		// g.Get("/admin/merchants/:merchantId/items", h.CreateMerchantItems)
	})
}
