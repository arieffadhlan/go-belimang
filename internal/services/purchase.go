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

func (s PurchaseService) GetNearbyMerchants(ctx context.Context, f entities.MerchantNearbyFilter) (map[string]any, error) {
	merchants, total, err := s.repository.GetNearbyMerchants(ctx, f)
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
				"name":             m.Name,
				"merchantId":       m.ID,
				"merchantCategory": m.Category,
				"imageUrl":         m.ImageURL,
				"location": map[string]float64{
					"lat":  m.Location.Lat,
					"long": m.Location.Lon,
				},
				"createdAt": m.CreatedAt.Format(time.RFC3339Nano),
			},
			"items": items,
		})
	}

	return map[string]any{
		"data": data,
		"meta": dto.Meta{Limit: f.Limit, Offset: f.Offset, Total: total},
	}, nil
}

func (s PurchaseService) CreateEstimate(ctx context.Context, req dto.EstimateReq) (dto.EstimateRes, error) {
	if err := ctx.Err(); err != nil {
		 return dto.EstimateRes{},err
	}

	startId := -1
	mercIDs := make([]string, 0, len(req.UserPurchase))
	itemIDs := make([]string, 0, len(req.UserPurchase)*3)

	for i, order := range req.UserPurchase {
		if order.IsStartingPoint {
			if startId >= 0 {
				 return dto.EstimateRes{}, utils.NewBadRequest("must have exactly one starting point")
			}
			startId = i
		}

		if _, err := uuid.Parse(order.MerchantID); err != nil {
				return dto.EstimateRes{}, utils.NewNotFound("merchant not found")
		} else {
				mercIDs = append(mercIDs, order.MerchantID)
		}

		for _, item := range order.OrderItems {
			if _, err := uuid.Parse(item.ItemID); err != nil {
				return dto.EstimateRes{}, utils.NewNotFound("mercItem not found")
			} else {
				itemIDs = append(itemIDs, item.ItemID)
			}
		}
	}

	if startId == -1 {
		 return dto.EstimateRes{}, utils.NewBadRequest("must have exactly one starting point")
	}

	merchants,err := s.repository.GetAllMerchantByIDs(ctx, mercIDs)
	if err != nil {
		 return dto.EstimateRes{}, utils.NewInternal("failed to get merchants")
	}
	if len(merchants) == 0 {
		 return dto.EstimateRes{}, utils.NewNotFound("merchant data not found")
	}

	mercItems,err := s.repository.GetAllMercItemByIDs(ctx, itemIDs)
	if err != nil {
		 return dto.EstimateRes{}, utils.NewInternal("failed to get mercItems")
	}
	if len(mercItems) == 0 {
		 return dto.EstimateRes{}, utils.NewNotFound("mercItem data not found")
	}

	merchantMap := make(map[string]entities.Merchant, len(merchants))
	mercItemMap := make(map[string]entities.MercItem, len(mercItems))
	for _, merchant := range merchants {
		merchantMap[merchant.ID] = merchant
	}
	for _, mercItem := range mercItems {
		mercItemMap[mercItem.ID] = mercItem
	}

	merchantPoints := make([]utils.Point, 0, len(req.UserPurchase)+1)
	for _, ord := range req.UserPurchase {
		merchant := merchantMap[ord.MerchantID]
		merchantPoints = append(merchantPoints, utils.Point{
			Lat: merchant.Location.Lat, 
			Lon: merchant.Location.Lon,
		})
	}

	merchantPoints = append(merchantPoints, utils.Point{
		Lat: req.UserLocation.Lat, 
		Lon: req.UserLocation.Lon,
	})

	maxDistance := 0.0
	usrLocPoint := merchantPoints[len(merchantPoints)-1]
	for _, p := range merchantPoints[:len(merchantPoints)-1] {
			if d := utils.Haversine(usrLocPoint, p); d > maxDistance {
				maxDistance = d
			}
	}

	if maxDistance > maxRadiusKm {
		return dto.EstimateRes{}, utils.NewBadRequest("distance too far")
	}

	totalDistance := utils.NearestNeighborTSP(startId, merchantPoints)

	totalPrice := 0
	orderItems := make([]entities.OrderItem, 0, len(itemIDs))
	for _, order := range req.UserPurchase {
		for _, orderItem := range order.OrderItems {
			item := mercItemMap[orderItem.ItemID]
			totalPrice += orderItem.ItemQuantity * item.Price
			orderItems = append(orderItems, entities.OrderItem{
				MerchantID:     order.MerchantID,
				MerchantItemID: orderItem.ItemID,
				Quantity:       orderItem.ItemQuantity,
			})
		}
	}

	tx,err := repository.BeginTx(ctx)
	if err != nil {
		 return dto.EstimateRes{}, err
	}
	defer tx.Rollback(ctx)

	estimateRq := entities.Estimate{UserID: req.UserID}
	estimateID, err := s.repository.CreateEstimateBatch(ctx, tx, estimateRq, orderItems)
	if err != nil {
		 return dto.EstimateRes{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		 return dto.EstimateRes{}, err
	}

	return dto.EstimateRes{
		TotalPrice:                   totalPrice,
		CalculatedEstimateId:         estimateID,
		EstimatedDeliveryTimeMinutes: int(totalDistance * kmToMinute),
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
		EstimateID: req.EstimateID,
	}

	tx,err := repository.BeginTx(ctx)
	if err != nil {
		 return dto.CreateOrderResponse{}, err
	}
	defer tx.Rollback(ctx)

	if err := s.repository.CreateOrderFromEsID(ctx, tx, order); err != nil {
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
