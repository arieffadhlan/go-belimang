package services

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/repository"
	"belimang/internal/utils"
	"context"
	"time"

	"github.com/google/uuid"
)

type PurchaseService struct {
	repository repository.PurchaseRepository
}

func NewPurchaseService(repository repository.PurchaseRepository) PurchaseService {
	return PurchaseService{repository: repository}
}

const (
	speedInKmh  = 40.0
	minPerHour  = 60.0
	kmToMinute  = 1 / speedInKmh * minPerHour
	maxRadiusKm = 3.0
)

func (s PurchaseService) GetNearbyMerchants(ctx context.Context, filter entities.MerchantNearbyFilter) (map[string]any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if filter.Lat == 0 && filter.Long == 0 {
		return nil, utils.NewBadRequest("invalid lat/long")
	}

	merchants, total, err := s.repository.GetNearbyMerchants(ctx, filter)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(merchants))
	for _, m := range merchants {
		items := make([]map[string]any, 0, len(m.Items))
		for _, it := range m.Items {
			items = append(items, map[string]any{
				"itemId":          it.ID,
				"name":            it.Name,
				"productCategory": it.Category,
				"price":           it.Price,
				"imageUrl":        it.ImageURL,
				"createdAt":       it.CreatedAt.Format(time.RFC3339Nano),
			})
		}

		data = append(data, map[string]any{
			"merchant": map[string]any{
				"merchantId":       m.ID,
				"name":             m.Name,
				"merchantCategory": m.Category,
				"imageUrl":         m.ImageURL,
				"location": map[string]float64{
					"lat":  m.Location.Lat,
					"long": m.Location.Long,
				},
				"createdAt": m.CreatedAt.Format(time.RFC3339Nano),
			},
			"items": items,
		})
	}

	resp := map[string]any{
		"data": data,
		"meta": dto.Meta{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}
	return resp, nil
}

func (s PurchaseService) CreateEstimate(ctx context.Context, req dto.EstimateRequest) (dto.EstimateResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.EstimateResponse{}, err
	}

	startId := -1
	mercIDs := make([]string, 0, len(req.UserPurchase))
	itemIDs := make([]string, 0, len(req.UserPurchase)*3)

	for i, order := range req.UserPurchase {
		if order.IsStartingPoint {
			if startId >= 0 {
				return dto.EstimateResponse{}, utils.NewBadRequest("must have exactly one starting point")
			}
			startId = i
		}

		if _, err := uuid.Parse(order.MerchantID); err != nil {
			return dto.EstimateResponse{}, utils.NewNotFound("merchant not found")
		} else {
			mercIDs = append(mercIDs, order.MerchantID)
		}

		for _, item := range order.OrderItems {
			if _, err := uuid.Parse(item.ItemID); err != nil {
				return dto.EstimateResponse{}, utils.NewNotFound("mercItem not found")
			} else {
				itemIDs = append(itemIDs, item.ItemID)
			}
		}
	}

	if startId == -1 {
		return dto.EstimateResponse{}, utils.NewBadRequest("must have exactly one starting point")
	}

	merchants, err := s.repository.GetAllMerchantByIDs(ctx, mercIDs)
	if err != nil {
		return dto.EstimateResponse{}, utils.NewInternal("failed to get merchants")
	}
	if len(merchants) == 0 {
		return dto.EstimateResponse{}, utils.NewNotFound("merchant data not found")
	}

	mercItems, err := s.repository.GetAllMercItemByIDs(ctx, itemIDs)
	if err != nil {
		return dto.EstimateResponse{}, utils.NewInternal("failed to get mercItems")
	}
	if len(mercItems) == 0 {
		return dto.EstimateResponse{}, utils.NewNotFound("mercItem data not found")
	}

	mercItemByID := make(map[string]entities.MerchantItem, len(mercItems))
	merchantByID := make(map[string]entities.Merchant, len(merchants))
	for _, merchant := range merchants {
		merchantByID[merchant.ID] = merchant
	}
	for _, mercItem := range mercItems {
		mercItemByID[mercItem.ID] = mercItem
	}

	points := make([]utils.Point, 0, len(req.UserPurchase)+1)
	for _, o := range req.UserPurchase {
		m := merchantByID[o.MerchantID]
		points = append(points, utils.Point{
			Lat: m.Location.Lat, Lon: m.Location.Long,
		})
	}

	points = append(points, utils.Point{
		Lat: req.UserLocation.Lat, Lon: req.UserLocation.Long,
	})

	maxDis := 0.0
	userPoint := points[len(points)-1]
	for _, point := range points[:len(points)-1] {
		d := utils.Haversine(userPoint, point)
		if d > maxDis {
			maxDis = d
		}
	}

	if maxDis > maxRadiusKm {
		return dto.EstimateResponse{}, utils.NewBadRequest("distance too far")
	}

	totalDist := utils.NearestNeighborTSP(startId, points)

	totalPrice := 0
	orderItems := make([]entities.OrderItem, 0, len(itemIDs))
	for _, order := range req.UserPurchase {
		for _, orderItem := range order.OrderItems {
			item := mercItemByID[orderItem.ItemID]
			totalPrice += orderItem.ItemQuantity * item.Price
			orderItems = append(orderItems, entities.OrderItem{
				ID:             uuid.New().String(),
				MerchantID:     order.MerchantID,
				MerchantItemID: orderItem.ItemID,
				Quantity:       orderItem.ItemQuantity,
			})
		}
	}

	tx, err := repository.BeginTx(ctx)
	if err != nil {
		return dto.EstimateResponse{}, err
	}
	defer tx.Rollback(ctx)

	est := entities.Estimate{
		ID:     uuid.New().String(),
		UserID: req.UserID,
	}
	if err := s.repository.CreateEstimateBatch(ctx, tx, est, orderItems); err != nil {
		return dto.EstimateResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.EstimateResponse{}, err
	}

	return dto.EstimateResponse{
		TotalPrice:                   totalPrice,
		CalculatedEstimateId:         est.ID,
		EstimatedDeliveryTimeMinutes: int(totalDist * kmToMinute),
	}, nil
}

func (s PurchaseService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return dto.CreateOrderResponse{}, err
	}

	_, err := s.repository.GetEstimateDataByID(ctx, req.EstimateID)
	if err != nil {
		return dto.CreateOrderResponse{}, utils.NewNotFound("estimate does not exist")
	}

	order := entities.Order{
		ID:         uuid.New().String(),
		EstimateID: req.EstimateID,
	}

	tx, err := repository.BeginTx(ctx)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}
	defer tx.Rollback(ctx)

	if err := s.repository.CreateOrderFromEsID(ctx, order); err != nil {
		return dto.CreateOrderResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.CreateOrderResponse{}, err
	}

	return dto.CreateOrderResponse{
		OrderID: order.ID,
	}, nil
}

func (s PurchaseService) GetAllOrder(ctx context.Context, filter entities.OrderFilter) ([]dto.OrderHistory, error) {
	if err := ctx.Err(); err != nil {
		return []dto.OrderHistory{}, err
	}

	return s.repository.GetAllOrder(ctx, filter)
}
