package services

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/repository"
	"belimang/internal/utils"
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

func (s MerchantService) CreateMerchantItems(ctx context.Context, req entities.MerchantItem) (dto.ItemMerchantResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.ItemMerchantResponse{}, err
	}

	tx, err := repository.BeginTx(ctx)
	if err != nil {
		return dto.ItemMerchantResponse{}, err
	}
	defer tx.Rollback(ctx)

	_, err = s.repository.GetMerchantById(ctx, tx, req.MerchantID)
	if err != nil {
		return dto.ItemMerchantResponse{}, utils.NewConflict("merchant do not exists")
	}

	res, err := s.repository.CreateItemMerchant(ctx, tx, req)
	if err != nil {
		return dto.ItemMerchantResponse{}, err
	}

	return dto.ItemMerchantResponse{
		ID: res.ID,
	}, nil
}

func (s MerchantService) GetAllMerchantItems(ctx context.Context, merchantId string, req entities.ItemFilter) (dto.ItemsResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.ItemsResponse{}, err
	}

	tx, err := repository.BeginTx(ctx)
	if err != nil {
		return dto.ItemsResponse{}, err
	}
	defer tx.Rollback(ctx)

	_, err = s.repository.GetMerchantById(ctx, tx, merchantId)
	if err != nil {
		return dto.ItemsResponse{}, utils.NewConflict("merchant do not exists")
	}

	res, err := s.repository.GetItems(ctx, tx, merchantId, req)

	if err != nil {
		return dto.ItemsResponse{}, err
	}

	return res, nil
}
