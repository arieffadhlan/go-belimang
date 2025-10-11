package route

import (
	"belimang/internal/handlers"
	"belimang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterPurchaseRoutes(r chi.Router, h handlers.PurchaseHandler) {
	r.Group(func(g chi.Router) {
		g.Use(middleware.Protected(false))

		g.Get("/users/orders", h.GetAllOrder)
		g.Get("/merchants/nearby/{lat},{lon}", h.GetNearbyMerchants)

		g.Post("/users/orders", h.CreateOrder)
		g.Post("/users/estimate", h.CreateEstimate)
	})
}
