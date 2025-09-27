package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantRepository struct {
	db *pgxpool.Pool
}

func NewMerchantRepository(db *pgxpool.Pool) MerchantRepository {
	return MerchantRepository{
		db: db,
	}
}
