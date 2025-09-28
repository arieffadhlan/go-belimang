package repository

import (
	"context"
	"belimang/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"belimang/internal/entities"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return AuthRepository{
		db: db,
	}
}

func (r AuthRepository) CreateUser(ctx context.Context, tx pgx.Tx, req entities.User) (entities.User, error) {
	if err := ctx.Err(); err != nil {
		return entities.User{}, err
	}

	query := `
		INSERT INTO users (
			email,
			username, 
			password, 
			is_admin
		)
		VALUES (
			$1, 
			$2, 
			$3, 
			$4 
		)
		RETURNING id, is_admin
	`

	usr := entities.User{}
	err := tx.QueryRow(
		ctx,
		query,
		req.Email,
		req.Username,
		req.Password,
		req.IsAdmin,
	).Scan(
		&usr.Id,
		&usr.IsAdmin,
	)

	if err != nil {
		return entities.User{}, utils.NewInternal("failed register account")
	}

	return usr, nil
}

func (r AuthRepository) GetUserByUsername(ctx context.Context, name string) (entities.User, error) {
	if err := ctx.Err(); err != nil {
		return entities.User{}, err
	}

	query := `SELECT id, password, is_admin FROM users WHERE username = $1 LIMIT 1`

	usr := entities.User{}
	err := r.db.QueryRow(ctx, query, name).Scan(&usr.Id, &usr.Password, &usr.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.User{}, utils.NewNotFound("users not found")
		} else {
			return entities.User{}, utils.NewInternal("failed get user")
		}
	}

	return usr, nil
}

func (r AuthRepository) GetUserByMailAddr(ctx context.Context, mail string) (entities.User, error) {
	if err := ctx.Err(); err != nil {
		return entities.User{}, err
	}

	query := `SELECT id, password, is_admin FROM users WHERE email = $1 LIMIT 1`

	usr := entities.User{}
	err := r.db.QueryRow(ctx, query, mail).Scan(&usr.Id, &usr.Password, &usr.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.User{}, utils.NewNotFound("users not found")
		} else {
			return entities.User{}, utils.NewInternal("failed get user")
		}
	}

	return usr, nil
}
