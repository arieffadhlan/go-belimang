package handlers

import "belimang/internal/services"

type MerchantHandler struct {
	service services.MerchantService
}

func NewMerchantHandler(service services.MerchantService) MerchantHandler {
	return MerchantHandler{
		service: service,
	}
}
