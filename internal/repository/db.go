package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func SetPool(p *pgxpool.Pool) {
	pool = p
}

func GetPool() *pgxpool.Pool {
	return pool
}

func BeginTx(ctx context.Context) (pgx.Tx, error) {
	return pool.BeginTx(ctx, pgx.TxOptions{})
}
