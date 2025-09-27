package route

import (
	"belimang/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterPurchaseRoutes(r chi.Router, h *handlers.PurchaseHandler) {
	r.Group(func(g chi.Router) {
		// g.Use(middleware.Protected)

		// g.Get("/merchants/nearby/:location", h.FindNearbyMerchant)
		// g.Get("/users/orders", h.GetEstimate)
		// g.Get("/users/orders", h.GetAllOrder)
		// g.Get("/users/orders", h.CreateOrder)
	})
}
