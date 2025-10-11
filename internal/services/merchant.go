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

func (s MerchantService) GetAllMerchant(ctx context.Context, req entities.MerchantFilter) (dto.MerchantResponse, error) {
	if err := ctx.Err(); err != nil {
		 return dto.MerchantResponse{}, err
	}

	return s.repository.GetAllMerchant(ctx, req)
}

func (s MerchantService) GetAllMercItem(ctx context.Context, req entities.MercItemFilter, merchantId string) (dto.MercItemResponse, error) {
	if err := ctx.Err(); err != nil {
		 return dto.MercItemResponse{}, err
	}

	_, err := s.repository.GetMerchantById(ctx, merchantId)
	if err != nil {
		 return dto.MercItemResponse{}, utils.NewNotFound("merchant does not exist")
	}

	return s.repository.GetAllMercItem(ctx, merchantId, req)
}

func NewMerchantService(repository repository.MerchantRepository) MerchantService {
	return MerchantService{
		repository: repository,
	}
}

func (s MerchantService) CreateMerchant(ctx context.Context, req entities.Merchant) (dto.CreateMerchantResponse, error) {
	if err := ctx.Err(); err != nil {
		 return dto.CreateMerchantResponse{}, err
	}

	tx,err := repository.BeginTx(ctx)
	if err != nil {
		 return dto.CreateMerchantResponse{}, err
	}
	defer tx.Rollback(ctx)

	merchants, err := s.repository.CreateMerchant(ctx, tx, req)
	if err != nil {
		 return dto.CreateMerchantResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		 return dto.CreateMerchantResponse{}, err
	}

	return dto.CreateMerchantResponse{
		ID: merchants.ID,
	}, nil
}

func (s MerchantService) CreateMercItem(ctx context.Context, req entities.MercItem) (dto.CreateMercItemResponse, error) {
	if err := ctx.Err(); err != nil {
		 return dto.CreateMercItemResponse{}, err
	}

	tx,err := repository.BeginTx(ctx)
	if err != nil {
		 return dto.CreateMercItemResponse{}, err
	}
	defer tx.Rollback(ctx)

	_, err = s.repository.GetMerchantById(ctx, req.MerchantID)
	if err != nil {
		 return dto.CreateMercItemResponse{}, utils.NewNotFound("merchant does not exist")
	}

	merchantItem, err := s.repository.CreateMercItem(ctx, tx, req)
	if err != nil {
		 return dto.CreateMercItemResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		 return dto.CreateMercItemResponse{}, err
	}

	return dto.CreateMercItemResponse{
		ID: merchantItem.ID,
	}, nil
}
