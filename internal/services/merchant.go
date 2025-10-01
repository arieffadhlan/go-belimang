package services

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/repository"
	"context"
)

type MerchantService struct {
	repository repository.MerchantRepository
}

func NewMerchantService(repository repository.MerchantRepository) MerchantService {
	return MerchantService{
		repository: repository,
	}
}

func (s MerchantService) CreateMerchant(ctx context.Context, req entities.Merchant) (dto.MerchantResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.MerchantResponse{}, err
	}

	tx, err := repository.BeginTx(ctx)
	if err != nil {
		return dto.MerchantResponse{}, err
	}
	defer tx.Rollback(ctx)

	res, err := s.repository.CreateMerchant(ctx, tx, req)
	if err != nil {
		return dto.MerchantResponse{}, err
	}

	return dto.MerchantResponse{
		ID: res.ID,
	}, nil
}

func (s MerchantService) GetMerchants(ctx context.Context, req entities.MerchantFilter) (dto.MerchantsResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.MerchantsResponse{}, err
	}

	tx, err := repository.BeginTx(ctx)
	if err != nil {
		return dto.MerchantsResponse{}, err
	}
	defer tx.Rollback(ctx)

	res, err := s.repository.GetMerchants(ctx, tx, req)
	if err != nil {
		return dto.MerchantsResponse{}, err
	}

	return res, nil
}
