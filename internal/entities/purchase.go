package entities

import "time"

type (
	Estimate struct {
		ID        string    `db:"id"`
		UserID    string    `db:"user_id"`
		CreatedAt time.Time `db:"created_at"`
	}

	Order struct {
		ID         string `db:"id"`
		EstimateID string `db:"estimate_id"`
	}

	OrderItem struct {
		ID             string `db:"id"`
		EstimateID     string `db:"estimate_id"`
		MerchantID     string `db:"merchant_id"`
		MerchantItemID string `db:"merchant_item_id"`
		Quantity       int    `db:"quantity"`
	}

	OrderFilter struct {
		Limit            int
		Name             string
		MerchantID       string
		MerchantCategory string
		UserID           string
		Offset           int
	}

	OrderDetail struct {
		OrderID           string    `db:"order_id" json:"orderId"`
		UserID            string    `db:"user_id" json:"userId"`
		MerchantID        string    `db:"merchant_id" json:"merchantId"`
		ItemID            string    `db:"item_id" json:"itemId"`
		MerchantName      string    `db:"merchant_name" json:"merchantName"`
		MerchantCategory  string    `db:"merchant_category" json:"merchantCategory"`
		MerchantImageURL  string    `db:"merchant_image_url" json:"merchantImageUrl"`
		ItemName          string    `db:"item_name" json:"itemName"`
		ItemCategory      string    `db:"item_category" json:"itemCategory"`
		ItemImageURL      string    `db:"item_image_url" json:"itemImageUrl"`
		MerchantLat       float64   `db:"merchant_lat" json:"merchantLat"`
		MerchantLong      float64   `db:"merchant_long" json:"merchantLong"`
		ItemPrice         int       `db:"item_price" json:"itemPrice"`
		Quantity          int       `db:"quantity" json:"quantity"`
		MerchantCreatedAt time.Time `db:"merchant_created_at" json:"merchantCreatedAt"`
		ItemCreatedAt     time.Time `db:"item_created_at" json:"itemCreatedAt"`
	}

	MerchantNearbyFilter struct {
		Name             string
		UserID           string
		MerchantID       string
		Limit, Offset    int
		MerchantCategory string
		Lat, Long        float64
	}

	MerchantWithItems struct {
		Merchant
		Items []MerchantItem
	}
)
