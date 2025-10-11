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

type PurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPurchaseRepository(db *pgxpool.Pool) PurchaseRepository {
	return PurchaseRepository{db: db}
}

func (r PurchaseRepository) GetNearbyMerchants(ctx context.Context, f entities.MerchantNearbyFilter) ([]entities.MerchantWithItems, int, error) {
	if err := ctx.Err(); err != nil {
		 return nil, 0, err
	}

	validEnums := map[string]bool{
		"SmallRestaurant": true, 
		"LargeRestaurant": true,
		"BoothKiosk": true, 
		"MediumRestaurant": true, 
		"ConvenienceStore": true,
		"MerchandiseRestaurant": true, 
	}

	if f.MerchantCategory != "" && !validEnums[f.MerchantCategory] {
		return []entities.MerchantWithItems{}, 0, nil
	}

	conds := []string{}
	args := []any{}
	i := 1

	// user location point
	args = append(args, f.Lon, f.Lat)
	lonIdx := i
	latIdx := i + 1
	i += 2

	// radius 3km filter
	conds = append(conds, fmt.Sprintf(
		"ST_DWithin(m.location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, 3000)",
		lonIdx, latIdx,
	))

	if f.MerchantID != "" {
		conds = append(conds, fmt.Sprintf("m.id = $%d::uuid", i))
		args = append(args, f.MerchantID)
		i++
	}

	if f.Name != "" {
		conds = append(conds, fmt.Sprintf(`(
			m.name ILIKE $%d OR 
			EXISTS (
				SELECT 1 FROM items it 
				WHERE it.merchant_id = m.id 
				AND it.name ILIKE $%d
			)
		)`, i, i))
		args = append(args, "%"+f.Name+"%")
		i++
	}

	if f.MerchantCategory != "" {
		conds = append(conds, fmt.Sprintf("m.category = $%d", i))
		args = append(args, f.MerchantCategory)
		i++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	limit := f.Limit
	if limit <= 0 {
		 limit = 5
	}

	offset := f.Offset
	if offset < 0 {
		 offset = 0
	}

	// count total merchants
	queryCount := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM merchants m
		%s
	`, where)

	var total int
	err := r.db.QueryRow(ctx, queryCount, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count merchants failed: %w", err)
	}

	if total == 0 {
		return []entities.MerchantWithItems{}, 0, nil
	}

	queryMerchants := fmt.Sprintf(`
		SELECT
			m.id::text,
			m.name,
			m.category,
			m.imageurl,
			ST_Y(m.location::geometry) AS lat,
			ST_X(m.location::geometry) AS lon,
			m.created_at,
			ST_Distance(m.location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography) as distance
		FROM merchants m
		%s
		ORDER BY distance
		LIMIT %d OFFSET %d
	`, lonIdx, latIdx, where, limit, offset)

	rows, err := r.db.Query(ctx, queryMerchants, args...)
	if err != nil {
		 return nil, 0, fmt.Errorf("query merchants failed: %w", err)
	}
	defer rows.Close()

	merchants := make([]entities.Merchant, 0)
	for rows.Next() {
		var m entities.Merchant
		var distance float64
		if err := rows.Scan(&m.ID, &m.Name, &m.Category, &m.ImageURL, &m.Location.Lat, &m.Location.Lon, &m.CreatedAt, &distance); err != nil {
			return nil, 0, fmt.Errorf("scan merchant failed: %w", err)
		}
		merchants = append(merchants, m)
	}
	if err := rows.Err(); err != nil {
		 return nil, 0, err
	}

	if len(merchants) == 0 {
		return []entities.MerchantWithItems{}, total, nil
	}

	ids := make([]string, 0, len(merchants))
	mmap := make(map[string]*entities.MerchantWithItems)
	for _, m := range merchants {
		ids = append(ids, m.ID)
		mmap[m.ID] = &entities.MerchantWithItems{Merchant: m, Items: []entities.MercItem{}}
	}

	itemQuery := `
		SELECT id::text, merchant_id::text, name, category, price, imageurl, created_at
		FROM items
		WHERE merchant_id = ANY($1)
		ORDER BY created_at DESC
	`

	itemRows, err := r.db.Query(ctx, itemQuery, ids)
	if err != nil {
		return nil, 0, fmt.Errorf("query items failed: %w", err)
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var it entities.MercItem
		if err := itemRows.Scan(&it.ID, &it.MerchantID, &it.Name, &it.Category, &it.Price, &it.ImageURL, &it.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan item failed: %w", err)
		}
		if m, ok := mmap[it.MerchantID]; ok {
			m.Items = append(m.Items, it)
		}
	}

	if err := itemRows.Err(); err != nil {
		return nil, 0, err
	}

	results := make([]entities.MerchantWithItems, 0, len(merchants))
	for _, m := range merchants {
		results = append(results, *mmap[m.ID])
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
		SELECT id, name, imageurl, category,
		       ST_X(location::geometry) AS lon,
		       ST_Y(location::geometry) AS lat,
		       created_at
		FROM merchants WHERE id = ANY($1)
	`, ids)
	if err != nil {
		 return nil, utils.NewInternal("query merchants failed")
	}
	defer rows.Close()

	merchants := make([]entities.Merchant, 0, len(ids))
	for rows.Next() {
		mrc := entities.Merchant{}
		err := rows.Scan(
			&mrc.ID,
			&mrc.Name,
			&mrc.ImageURL,
			&mrc.Category,
			&mrc.Location.Lon,
			&mrc.Location.Lat,
			&mrc.CreatedAt,
		) 

		if err != nil {
			 return nil, utils.NewInternal("failed to scan merchant row")
		}

		merchants = append(merchants, mrc)
	}

	if err := rows.Err(); err != nil {
		 return nil, utils.NewInternal("error iterating merchant rows")
	}

	return merchants, nil
}

func (r PurchaseRepository) GetAllMercItemByIDs(ctx context.Context, ids []string) ([]entities.MercItem, error) {
	if err := ctx.Err(); err != nil {
		 return nil, err
	}

	if len(ids) == 0 {
		 return []entities.MercItem{}, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, merchant_id, name, price, imageurl, category, created_at
		FROM items
		WHERE id = ANY($1)
	`, ids)
	if err != nil {
		 return nil, utils.NewInternal("query mercItems failed")
	}
	defer rows.Close()

	items := make([]entities.MercItem, 0, len(ids))
	for rows.Next() {
		itm := entities.MercItem{}
		err := rows.Scan(
			&itm.ID,
			&itm.MerchantID,
			&itm.Name,
			&itm.Price,
			&itm.ImageURL,
			&itm.Category,
			&itm.CreatedAt,
		); 

		if err != nil {
			 return nil, utils.NewInternal("failed to scan merchant items row")
		}
		
		items = append(items, itm)
	}

	if err := rows.Err(); err != nil {
		 return nil, utils.NewInternal("error iterating merchant items rows")
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

func (r PurchaseRepository) CreateEstimateBatch(ctx context.Context, tx pgx.Tx, est entities.Estimate, items []entities.OrderItem) (string, error) {
	if err := ctx.Err(); err != nil {
		 return "", err
	}

	if len(items) == 0 {
		 return "", nil
	}

	var estimateID string
	err := tx.QueryRow(ctx, `INSERT INTO estimates (user_id) VALUES ($1) RETURNING id`, est.UserID).Scan(&estimateID)
	if err != nil {
		 return "", utils.NewInternal("failed to insert estimate")
	}

	batch := &pgx.Batch{}
	for _, it := range items {
		batch.Queue(`
			INSERT INTO orders_items (
				estimate_id, 
				merchant_id, 
				merchant_item_id, 
				quantity
			)
			VALUES ($1, $2, $3, $4)
		`, 
			estimateID, 
			it.MerchantID, 
			it.MerchantItemID, 
			it.Quantity,
		)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()
	if err := br.Close(); err != nil {
		return "", utils.NewInternal("failed to batch insert order items")
	}

	return estimateID, nil
}

func (r PurchaseRepository) CreateOrderFromEsID(ctx context.Context, tx pgx.Tx, order entities.Order) error {
	if err := ctx.Err(); err != nil {
		 return err
	}

	_, err := tx.Exec(ctx, `INSERT INTO orders (estimate_id) VALUES ($1)`, order.EstimateID)
	if err != nil {
		return err
	}

	return nil
}

type OrderGroup struct {
	Order *dto.OrderHistory
	Group map[string]*dto.OrderHistoryMerchant
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

	query := fmt.Sprintf(`
		SELECT
			order_id,
			user_id,
			merchant_id,
			merchant_name,
			merchant_category,
			merchant_imageurl,
			merchant_lat,
			merchant_lon,
			merchant_created_at,
			item_id,
			item_name,
			item_category,
			item_imageurl,
			item_price,
			quantity,
			item_created_at
		FROM order_history_view
		WHERE %s
		ORDER BY order_id DESC
		LIMIT %d OFFSET %d
	`, strings.Join(conditions, " AND "), limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		 return nil, utils.NewInternal("failed to query order history view")
	}
	defer rows.Close()

	orderMap := make(map[string]*OrderGroup, 512)
	merchantCache := make(map[string]dto.Merchant, 872)

	for rows.Next() {
		ord := entities.OrderDetail{}
		err := rows.Scan(
			&ord.OrderID,
			&ord.UserID,
			&ord.MerchantID,
			&ord.MerchantName,
			&ord.MerchantCategory,
			&ord.MerchantImageURL,
			&ord.MerchantLat,
			&ord.MerchantLon,
			&ord.MerchantCreatedAt,
			&ord.ItemID,
			&ord.ItemName,
			&ord.ItemCategory,
			&ord.ItemImageURL,
			&ord.ItemPrice,
			&ord.Quantity,
			&ord.ItemCreatedAt,
		)

		if err != nil {
			 return nil, utils.NewInternal("failed to scan order history row")
		}

		mrc,found := merchantCache[ord.MerchantID]
		if !found {
			mrc = dto.Merchant{
				ID:        ord.MerchantID,
				Name:      ord.MerchantName,
				CreatedAt: ord.MerchantCreatedAt,
				Category:  ord.MerchantCategory,
				ImageURL:  ord.MerchantImageURL,
				Location: dto.Location{
					Lat: ord.MerchantLat,
					Lon: ord.MerchantLon,
				},
			}
			
			merchantCache[ord.MerchantID] = mrc
		}

		grp,ok := orderMap[ord.OrderID]
		if !ok {
			grp = &OrderGroup{
				Group:make(map[string]*dto.OrderHistoryMerchant, 64),
				Order:&dto.OrderHistory{
					OrderID:      ord.OrderID,
					OrderHistory: make([]dto.OrderHistoryMerchant, 0, 8),
				},
			}

			orderMap[ord.OrderID] = grp
		}

		merchantGroup, ok := grp.Group[ord.MerchantID]
		if !ok {
			newGroup := dto.OrderHistoryMerchant{
				Merchant: mrc,
				Items:    make([]dto.OrderItem, 0, 16),
			}
			grp.Order.OrderHistory = append(grp.Order.OrderHistory, newGroup)
			merchantGroup = &grp.Order.OrderHistory[len(grp.Order.OrderHistory)-1]
			grp.Group[ord.MerchantID] = merchantGroup
		}

		merchantGroup.Items = append(merchantGroup.Items, dto.OrderItem{
			ItemID:          ord.ItemID,
			Name:            ord.ItemName,
			Quantity:        ord.Quantity,
			ImageURL:        ord.ItemImageURL,
			ProductCategory: ord.ItemCategory,
			Price:           ord.ItemPrice,
			CreatedAt:       ord.ItemCreatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		 return nil, utils.NewInternal("error iterating order history rows")
	}

	result := make([]dto.OrderHistory, 0, len(orderMap))
	for _, v := range orderMap {
		 result = append(result, *v.Order)
	}

	return result, nil
}

