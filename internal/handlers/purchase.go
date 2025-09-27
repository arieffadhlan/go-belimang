package handlers

import "belimang/internal/services"

type PurchaseHandler struct {
	service services.PurchaseService
}

func NewPurchaseHandler(service services.PurchaseService) PurchaseHandler {
	return PurchaseHandler{
		service: service,
	}
}
