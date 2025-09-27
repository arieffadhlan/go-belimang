package services

import "belimang/internal/repository"

type PurchaseService struct {
	repository repository.PurchaseRepository
}

func NewPurchaseService(repository repository.PurchaseRepository) PurchaseService {
	return PurchaseService{
		repository: repository,
	}
}
