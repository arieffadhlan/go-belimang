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
	return MerchantRepository{db: db}
}

func (r MerchantRepository) GetAllMerchant(ctx context.Context, filter entities.MerchantFilter) (dto.MerchantResponse, error) {
	if err := ctx.Err(); err != nil {
		 return dto.MerchantResponse{}, err
	}

	conditions := []string{"1=1"}
	args := []any{}
	i := 1

	if filter.MerchantCategory != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", i))
		args = append(args, filter.MerchantCategory)
		i++
	}

	if filter.MerchantID != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", i))
		args = append(args, filter.MerchantID)
		i++
	}

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", i))
		args = append(args, "%"+filter.Name+"%")
		i++
	}

	order := "DESC"
	if filter.CreatedAt == "asc" || filter.CreatedAt == "desc" {
		 order = strings.ToUpper(filter.CreatedAt)
	}

	limit, offset := filter.Limit, filter.Offset
	if limit <= 0 {
		 limit = 5
	}
	
	if offset < 0 {
		 offset = 0
	}

	query := fmt.Sprintf(`
		SELECT 
			id, name, imageurl, category,
			ST_X(location::geometry) as lon,
			ST_Y(location::geometry) as lat,
			created_at, Count(*) OVER() AS total
		FROM merchants
		WHERE %s
		ORDER BY created_at %s
		LIMIT %d OFFSET %d
	`, strings.Join(conditions, " AND "), order, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		 return dto.MerchantResponse{}, utils.NewInternal("failed to query merchants")
	}
	defer rows.Close()

	var total int
	var merchants []dto.Merchant

	for rows.Next() {
		cur := dto.Merchant{}
		err := rows.Scan(
			&cur.ID,
			&cur.Name,
			&cur.ImageURL,
			&cur.Category,
			&cur.Location.Lon, 
			&cur.Location.Lat,
			&cur.CreatedAt,
			&total,
		)

		if err != nil {
			 return dto.MerchantResponse{}, utils.NewInternal("failed to scan merchant")
		}

		merchants = append(merchants, cur)
	}

	if merchants == nil {
		 merchants = make([]dto.Merchant, 0)
	}

	return dto.MerchantResponse{
		Data: merchants,
		Meta: dto.Meta{Total: total, Limit: limit, Offset: offset},
	}, nil
}

func (r MerchantRepository) GetAllMercItem(ctx context.Context, merchantId string, filter entities.MercItemFilter) (dto.MercItemResponse, error) {
	if err := ctx.Err(); err != nil {
		 return dto.MercItemResponse{}, err
	}

	conditions := []string{"merchant_id = $1"}
	args := []any{merchantId}
	i := 2

	if filter.ProductCategory != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", i))
		args = append(args, filter.ProductCategory)
		i++
	}

	if filter.ItemID != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", i))
		args = append(args, filter.ItemID)
		i++
	}

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", i))
		args = append(args, "%"+filter.Name+"%")
		i++
	}

	order := "DESC"
	if filter.CreatedAt == "asc" || filter.CreatedAt == "desc" {
		 order = strings.ToUpper(filter.CreatedAt)
	}

	limit, offset := filter.Limit, filter.Offset
	if limit <= 0 {
		 limit = 5
	}

	if offset < 0 {
		 offset = 0
	}
	
	query := fmt.Sprintf(`
		SELECT 
			id, name, price, imageurl, category, created_at, 
			COUNT(*) OVER() AS total
		FROM items 
		WHERE %s
		ORDER BY created_at %s
		LIMIT %d OFFSET %d
	`, strings.Join(conditions, " AND "), order, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		 return dto.MercItemResponse{}, utils.NewInternal("failed to query items")
	}
	defer rows.Close()

	var total int
	var items []dto.MercItem
	
	for rows.Next() {
		var item dto.MercItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.ImageURL,
			&item.Category,
			&item.CreateAt,
			&total,
		)

		if err != nil {
			 return dto.MercItemResponse{}, utils.NewInternal("failed to scan merchant item")
		}

		items = append(items, item)
	}

	if items == nil {
		 items = make([]dto.MercItem, 0)
	}

	return dto.MercItemResponse{
		Data: items,
		Meta: dto.Meta{Total: total, Limit: limit, Offset: offset},
	}, nil
}

func (r MerchantRepository) GetMerchantById(ctx context.Context, merchantId string) (entities.Merchant, error) {
	if err := ctx.Err(); err != nil {
		 return entities.Merchant{}, err
	}

	query := `SELECT id FROM merchants WHERE id = $1`

	var m entities.Merchant
	err := r.db.QueryRow(ctx, query, merchantId).Scan(&m.ID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.Merchant{}, utils.NewNotFound("merchants not found")
		} else {
			return entities.Merchant{}, utils.NewInternal("failed get merchant")
		}
	}

	return m, nil
}

func (r MerchantRepository) CreateMerchant(ctx context.Context, tx pgx.Tx, req entities.Merchant) (entities.Merchant, error) {
	if err := ctx.Err(); err != nil {
		 return entities.Merchant{}, err
	}

	query := `
		INSERT INTO merchants (
			name, 
			imageurl, 
			category, 
			location
		)
		VALUES (
			$1, $2, $3,
			ST_SetSRID(ST_MakePoint($4, $5), 4326)::GEOGRAPHY
		)
		RETURNING id
	`

	res := entities.Merchant{}
	err := tx.QueryRow(ctx, query,
		req.Name,
		req.ImageURL,
		req.Category,
		req.Location.Lon, 
		req.Location.Lat,
	).Scan(&res.ID)

	if err != nil {
		 return entities.Merchant{}, err
	}

	return res, nil
}

func (r MerchantRepository) CreateMercItem(ctx context.Context, tx pgx.Tx, req entities.MercItem) (entities.MercItem, error) {
	if err := ctx.Err(); err != nil {
		 return entities.MercItem{}, err
	}

	query := `
		INSERT INTO items (merchant_id, name, price, imageurl, category)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	res := entities.MercItem{}
	err := tx.QueryRow(ctx, query, req.MerchantID, req.Name, req.Price, req.ImageURL, req.Category).Scan(&res.ID)

	if err != nil {
		 return entities.MercItem{}, utils.NewInternal("failed create merchant item")
	}

	return res, nil
}
