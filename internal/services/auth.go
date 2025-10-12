package services

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/repository"
	"belimang/internal/utils"
	"context"
)

type AuthService struct {
	repository repository.AuthRepository
	hashWkPool *HashingWorkerPool
}

func NewAuthService(repository repository.AuthRepository, hashWkPool *HashingWorkerPool) AuthService {
	return AuthService{
		repository: repository,
		hashWkPool: hashWkPool,
	}
}

func (s AuthService) SignUp(ctx context.Context, req entities.User) (dto.AuthResponse, error) {
	if err := ctx.Err(); err != nil {
		 return dto.AuthResponse{}, err
	}

	hash, err := s.hashWkPool.HashPasswordAsync(ctx, req.Password)
	if err != nil {
		 return dto.AuthResponse{}, err
	}

	req.Password = hash

	tx,err := repository.BeginTx(ctx)
	if err != nil {
		 return dto.AuthResponse{}, err
	}
	defer tx.Rollback(ctx)

	_, err = s.repository.GetUserByUsername(ctx, req.Username)
	if err == nil {
		 return dto.AuthResponse{}, utils.NewConflict("account already exists")
	}

	e, err := s.repository.GetUserByMailAddr(ctx, req.Email)
	if err == nil && e.Id != "" && e.IsAdmin != req.IsAdmin {
		 return dto.AuthResponse{}, utils.NewConflict("account already exists")
	}

	u, err := s.repository.CreateUser(ctx, tx, req)
	if err != nil {
		 return dto.AuthResponse{}, err
	}

	t, err := GenerateToken(u)
	if err != nil {
		 return dto.AuthResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		 return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token: t,
	}, nil
}

func (s AuthService) SignIn(ctx context.Context, req entities.User) (dto.AuthResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.AuthResponse{}, err
	}

	u, err := s.repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if !ComparePassword(req.Password, u.Password) {
		return dto.AuthResponse{}, utils.NewBadRequest("invalid credentials")
	}

	t, err := GenerateToken(u)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token: t,
	}, nil
}
