package services

import "belimang/internal/repository"

type MerchantService struct {
	repository repository.MerchantRepository
}

func NewMerchantService(repository repository.MerchantRepository) MerchantService {
	return MerchantService{
		repository: repository,
	}
}
