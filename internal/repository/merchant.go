package repository

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/utils"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
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

func (r MerchantRepository) CreateMerchant(ctx context.Context, tx pgx.Tx, req entities.Merchant) (entities.Merchant, error) {
	if err := ctx.Err(); err != nil {
		return entities.Merchant{}, err
	}

	query := `
		INSERT INTO merchants (name, merchant_category, image_url, location)
		VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326)::GEOGRAPHY)
		RETURNING merchant_id, name, merchant_category, image_url,
		          ST_Y(location::geometry) AS lat,
		          ST_X(location::geometry) AS long,
		          created_at
	`

	merchant := entities.Merchant{}
	err := tx.QueryRow(
		ctx,
		query,
		req.Name,
		req.MerchantCategory,
		req.ImageURL,
		req.Location.Long,
		req.Location.Lat,
	).Scan(&merchant.ID)

	if err != nil {
		return entities.Merchant{}, utils.NewInternal("failed create merchant")
	}

	return merchant, nil
}

func (r MerchantRepository) GetMerchants(ctx context.Context, tx pgx.Tx, filter entities.MerchantFilter) (dto.MerchantsResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.MerchantsResponse{}, err
	}

	args := []interface{}{}
	conditions := []string{"1=1"} // always true so we can append conditions dynamically
	argIndex := 1                 // pgx uses $1, $2... placeholders

	if filter.MerchantID != "" {
		conditions = append(conditions, fmt.Sprintf("merchant_id = $%d", argIndex))
		args = append(args, filter.MerchantID)
		argIndex++
	}

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE $%d", argIndex))
		args = append(args, "%"+strings.ToLower(filter.Name)+"%")
		argIndex++
	}

	if filter.MerchantCategory != "" {
		conditions = append(conditions, fmt.Sprintf("merchant_category = $%d", argIndex))
		args = append(args, filter.MerchantCategory)
		argIndex++
	}

	// Base query
	baseQuery := `
		SELECT 
			merchant_id,
			name,
			merchant_category,
			image_url,
			ST_Y(location::geometry) as lat,
			ST_X(location::geometry) as long,
			created_at
		FROM merchants
		WHERE ` + strings.Join(conditions, " AND ")

	// Sorting
	if filter.SortCreatedAt == "asc" || filter.SortCreatedAt == "desc" {
		baseQuery += fmt.Sprintf(" ORDER BY created_at %s", filter.SortCreatedAt)
	}

	// Pagination
	limit := 5
	offset := 0
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Offset >= 0 {
		offset = filter.Offset
	}
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Execute query
	rows, err := tx.Query(ctx, baseQuery, args...)
	if err != nil {
		return dto.MerchantsResponse{}, utils.NewInternal("failed to query merchants")
	}
	defer rows.Close()

	var merchants []dto.Merchant
	for rows.Next() {
		var m dto.Merchant
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.MerchantCategory,
			&m.ImageURL,
			&m.Location.Lat,
			&m.Location.Long,
			&m.CreatedAt,
		)
		if err != nil {
			return dto.MerchantsResponse{}, utils.NewInternal("failed to scan merchant")
		}
		merchants = append(merchants, m)
	}

	// Count query for meta
	countQuery := `
		SELECT COUNT(*)
		FROM merchants
		WHERE ` + strings.Join(conditions, " AND ")

	var total int
	err = tx.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return dto.MerchantsResponse{}, utils.NewInternal("failed to count merchants")
	}

	return dto.MerchantsResponse{
		Data: merchants,
		Meta: dto.Meta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}, nil
}
