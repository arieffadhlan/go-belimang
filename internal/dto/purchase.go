package dto

import "time"

type (
	EstimateReq struct {
		UserID       string          `json:"_"`
		UserPurchase []EstimateOrder `json:"orders" validate:"required,dive"`
		UserLocation Location        `json:"userLocation" validate:"required"`
	}

	EstimateRes struct {
		TotalPrice                   int    `json:"totalPrice"`
		CalculatedEstimateId         string `json:"calculatedEstimateId"`
		EstimatedDeliveryTimeMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	}

	EstimateOrder struct {
		IsStartingPoint bool                `json:"isStartingPoint" validate:"boolean"`
		OrderItems      []EstimateOrderItem `json:"items" validate:"required,dive"`
		MerchantID      string              `json:"merchantId" validate:"required"`
	}

	EstimateOrderItem struct {
		ItemID       string `json:"itemId"   validate:"required"`
		ItemQuantity int    `json:"quantity" validate:"required,gt=0"`
	}

	CreateOrderRequest struct {
		EstimateID string `json:"calculatedEstimateId" validate:"required"`
	}

	CreateOrderResponse struct {
		OrderID string `json:"orderId"`
	}

	OrderHistory struct {
		OrderID      string                 `json:"orderId"`
		OrderHistory []OrderHistoryMerchant `json:"orders"`
	}

	OrderHistoryMerchant struct {
		Merchant Merchant    `json:"merchant"`
		Items    []OrderItem `json:"items"`
	}

	OrderItem struct {
		ItemID          string    `json:"itemId"`
		ProductCategory string    `json:"productCategory"`
		Name            string    `json:"name"`
		Quantity        int       `json:"quantity"`
		ImageURL        string    `json:"imageUrl"`
		Price           int       `json:"price"`
		CreatedAt       time.Time `json:"createdAt"`
	}
)
