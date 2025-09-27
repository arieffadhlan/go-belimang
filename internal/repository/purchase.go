package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPurchaseRepository(db *pgxpool.Pool) PurchaseRepository {
	return PurchaseRepository{
		db: db,
	}
}
