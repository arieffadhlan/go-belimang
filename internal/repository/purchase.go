package repository

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/utils"
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPurchaseRepository(db *pgxpool.Pool) PurchaseRepository {
	return PurchaseRepository{db: db}
}

func (r PurchaseRepository) GetNearbyMerchants(
	ctx context.Context,
	f entities.MerchantNearbyFilter,
) ([]entities.MerchantWithItems, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	// Build dynamic WHERE
	conds := []string{"TRUE"}
	args := []any{}
	i := 1

	if f.MerchantID != "" {
		conds = append(conds, fmt.Sprintf("m.id::text = $%d", i))
		args = append(args, f.MerchantID)
		i++
	}
	if f.MerchantCategory != "" {
		validEnums := map[string]bool{
			"SmallRestaurant": true, "MediumRestaurant": true, "LargeRestaurant": true,
			"MerchandiseRestaurant": true, "BoothKiosk": true, "ConvenienceStore": true,
		}
		if !validEnums[f.MerchantCategory] {
			return []entities.MerchantWithItems{}, 0, nil
		}
		conds = append(conds, fmt.Sprintf("m.category = $%d", i))
		args = append(args, f.MerchantCategory)
		i++
	}
	if f.Name != "" {
		conds = append(conds, fmt.Sprintf(
			`(m.name ILIKE $%d OR EXISTS (
				SELECT 1 FROM items it
				WHERE it.merchant_id::uuid = m.id
				  AND it.name ILIKE $%d
			))`, i, i))
		args = append(args, "%"+f.Name+"%")
		i++
	}

	where := "WHERE " + strings.Join(conds, " AND ")

	query := fmt.Sprintf(`
		SELECT 
			m.id::text, m.name, m.category, m.image_url,
			ST_Y(m.location::geometry) AS lat,
			ST_X(m.location::geometry) AS lon,
			m.created_at
		FROM merchants m
		%s
	`, where)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, utils.NewInternal("failed to query merchants")
	}
	defer rows.Close()

	// Ambil semua merchant
	all := []entities.Merchant{}
	for rows.Next() {
		var m entities.Merchant
		if err := rows.Scan(
			&m.ID, &m.Name, &m.Category, &m.ImageURL,
			&m.Location.Lat, &m.Location.Long, &m.CreatedAt,
		); err != nil {
			return nil, 0, utils.NewInternal("failed to scan merchant row")
		}
		all = append(all, m)
	}

	user := utils.Point{Lat: f.Lat, Lon: f.Long}
	const maxRadiusKm = 3.0

	type mwrap struct {
		entities.Merchant
		dist float64
	}
	filtered := []mwrap{}

	for _, m := range all {
		d := utils.HaversineKm(user, utils.Point{Lat: m.Location.Lat, Lon: m.Location.Long})
		if d <= maxRadiusKm {
			filtered = append(filtered, mwrap{Merchant: m, dist: d})
		}
	}

	if len(filtered) == 0 {
		return []entities.MerchantWithItems{}, 0, nil
	}

	// 🔥 Urutkan berdasar jarak langsung (bukan TSP)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].dist < filtered[j].dist
	})

	// Ambil 5 terdekat
	if len(filtered) > 5 {
		filtered = filtered[:5]
	}
	total := len(filtered)

	// Query item per merchant
	ids := make([]string, 0, total)
	mmap := make(map[string]*entities.MerchantWithItems, total)
	for _, mw := range filtered {
		ids = append(ids, mw.ID)
		mmap[mw.ID] = &entities.MerchantWithItems{
			Merchant: mw.Merchant,
			Items:    []entities.MerchantItem{},
		}
	}

	itemRows, err := r.db.Query(ctx, `
		SELECT id, merchant_id, name, category, price, image_url, created_at
		FROM items
		WHERE merchant_id = ANY($1)
		ORDER BY created_at DESC
	`, ids)
	if err != nil {
		return nil, 0, utils.NewInternal("failed to query items")
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var it entities.MerchantItem
		if err := itemRows.Scan(
			&it.ID, &it.MerchantID, &it.Name, &it.Category,
			&it.Price, &it.ImageURL, &it.CreatedAt,
		); err != nil {
			return nil, 0, utils.NewInternal("failed to scan item row")
		}
		if m, ok := mmap[it.MerchantID]; ok {
			m.Items = append(m.Items, it)
		}
	}

	results := make([]entities.MerchantWithItems, 0, total)
	for _, mw := range filtered {
		results = append(results, *mmap[mw.ID])
	}

	return results, total, nil
}

func (r PurchaseRepository) GetAllMerchantByIDs(ctx context.Context, ids []string) ([]entities.Merchant, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []entities.Merchant{}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, name, image_url, category,
		       ST_X(location::geometry) AS lon,
		       ST_Y(location::geometry) AS lat,
		       created_at
		FROM merchants WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return nil, utils.NewInternal("query merchants failed")
	}
	defer rows.Close()

	merchants, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.Merchant, error) {
		mrc := entities.Merchant{}
		err := row.Scan(&mrc.ID, &mrc.Name, &mrc.ImageURL, &mrc.Category, &mrc.Location.Long, &mrc.Location.Lat, &mrc.CreatedAt)
		return mrc, err
	})

	if err != nil {
		return nil, utils.NewInternal("scan merchants failed")
	}

	return merchants, nil
}

func (r PurchaseRepository) GetAllMercItemByIDs(ctx context.Context, ids []string) ([]entities.MerchantItem, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []entities.MerchantItem{}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, merchant_id, name, price, image_url, category, created_at
		FROM items
		WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return nil, utils.NewInternal("query mercItems failed")
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.MerchantItem, error) {
		itm := entities.MerchantItem{}
		err := row.Scan(&itm.ID, &itm.MerchantID, &itm.Name, &itm.Price, &itm.ImageURL, &itm.Category, &itm.CreatedAt)
		return itm, err
	})

	if err != nil {
		return nil, utils.NewInternal("scan mercItems failed")
	}

	return items, nil
}

func (r PurchaseRepository) GetEstimateDataByID(ctx context.Context, id string) (entities.Estimate, error) {
	if err := ctx.Err(); err != nil {
		return entities.Estimate{}, err
	}

	row := r.db.QueryRow(ctx, `SELECT id, user_id, created_at FROM estimates WHERE id = $1`, id)

	est := entities.Estimate{}
	err := row.Scan(&est.ID, &est.UserID, &est.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.Estimate{}, utils.NewNotFound("estimates not found")
		} else {
			return entities.Estimate{}, utils.NewInternal("failed get estimate")
		}
	}

	return est, nil
}

func (r PurchaseRepository) CreateEstimateBatch(ctx context.Context, tx pgx.Tx, est entities.Estimate, items []entities.OrderItem) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	var estimateID string
	err := tx.QueryRow(ctx, `INSERT INTO estimates (id, user_id) VALUES ($1, $2) RETURNING id`, est.ID, est.UserID).Scan(&estimateID)
	if err != nil {
		return utils.NewInternal("failed to insert estimate")
	}

	batch := &pgx.Batch{}
	for _, it := range items {
		batch.Queue(`
			INSERT INTO orders_items (
				id,
				estimate_id, 
				merchant_id, 
				merchant_item_id, 
				quantity
			)
			VALUES ($1, $2, $3, $4, $5)
		`, it.ID, estimateID, it.MerchantID, it.MerchantItemID, it.Quantity)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()
	if err := br.Close(); err != nil {
		return utils.NewInternal("failed to batch insert order items")
	}

	return nil
}

func (r PurchaseRepository) CreateOrderFromEsID(ctx context.Context, order entities.Order) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.Exec(ctx, `INSERT INTO orders (id, estimate_id) VALUES ($1, $2)`, order.ID, order.EstimateID)
	if err != nil {
		return err
	}

	return nil
}

func (r PurchaseRepository) GetAllOrder(ctx context.Context, filter entities.OrderFilter) ([]dto.OrderHistory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	limit, offset := filter.Limit, filter.Offset
	if limit <= 0 {
		limit = 5
	}
	if offset < 0 {
		offset = 0
	}

	conditions := []string{"user_id = $1"}
	args := []any{filter.UserID}
	i := 2

	if filter.MerchantCategory != "" {
		conditions = append(conditions, fmt.Sprintf("merchant_category = $%d", i))
		args = append(args, filter.MerchantCategory)
		i++
	}

	if filter.MerchantID != "" {
		conditions = append(conditions, fmt.Sprintf("merchant_id = $%d", i))
		args = append(args, filter.MerchantID)
		i++
	}

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("(merchant_name ILIKE $%d OR item_name ILIKE $%d)", i, i))
		args = append(args, "%"+filter.Name+"%")
		i++
	}

	where := strings.Join(conditions, " AND ")
	query := fmt.Sprintf(`
		SELECT
			order_id,
			user_id,
			merchant_id,
			merchant_name,
			merchant_category,
			merchant_image_url,
			merchant_lat,
			merchant_long,
			merchant_created_at,
			item_id,
			item_name,
			item_category,
			item_price,
			quantity,
			item_image_url,
			item_created_at
		FROM order_history_view
		WHERE %s
		ORDER BY order_id DESC
		LIMIT %d OFFSET %d
	`, where, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, utils.NewInternal("failed to query order history view")
	}
	defer rows.Close()

	type key struct {
		OrderID, MerchantID string
	}

	itemsMap := make(map[key][]dto.OrderItemDTO, 10)
	orderMap := make(map[string][]dto.OrderHistoryMerchant, 10)
	merchantCache := make(map[string]dto.Merchant, 60)

	for rows.Next() {
		var d entities.OrderDetail
		if err := rows.Scan(
			&d.OrderID,
			&d.UserID,
			&d.MerchantID,
			&d.MerchantName,
			&d.MerchantCategory,
			&d.MerchantImageURL,
			&d.MerchantLat,
			&d.MerchantLong,
			&d.MerchantCreatedAt,
			&d.ItemID,
			&d.ItemName,
			&d.ItemCategory,
			&d.ItemPrice,
			&d.Quantity,
			&d.ItemImageURL,
			&d.ItemCreatedAt,
		); err != nil {
			return nil, utils.NewInternal("failed to scan order history row")
		}

		if _, exists := merchantCache[d.MerchantID]; !exists {
			merchantCache[d.MerchantID] = dto.Merchant{
				ID:        d.MerchantID,
				Name:      d.MerchantName,
				Category:  d.MerchantCategory,
				ImageURL:  d.MerchantImageURL,
				Location:  dto.Location{Lat: d.MerchantLat, Long: d.MerchantLong},
				CreatedAt: d.MerchantCreatedAt,
			}
		}

		k := key{OrderID: d.OrderID, MerchantID: d.MerchantID}
		itemsMap[k] = append(itemsMap[k], dto.OrderItemDTO{
			ItemID:          d.ItemID,
			ProductCategory: d.ItemCategory,
			Name:            d.ItemName,
			Quantity:        d.Quantity,
			ImageURL:        d.ItemImageURL,
			Price:           d.ItemPrice,
			CreatedAt:       d.ItemCreatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewInternal("error iterating order history rows")
	}

	for k, items := range itemsMap {
		orderMap[k.OrderID] = append(orderMap[k.OrderID], dto.OrderHistoryMerchant{
			Merchant: merchantCache[k.MerchantID],
			Items:    items,
		})
	}

	result := make([]dto.OrderHistory, 0, len(orderMap))
	for orderID, orders := range orderMap {
		result = append(result, dto.OrderHistory{
			OrderID: orderID,
			Orders:  orders,
		})
	}

	return result, nil
}
