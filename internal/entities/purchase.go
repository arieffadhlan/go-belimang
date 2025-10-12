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
		UserID            string    `db:"user_id" json:"userId"`
		ItemID            string    `db:"item_id" json:"itemId"`
		MerchantID        string    `db:"merchant_id" json:"merchantId"`
		MerchantImageURL  string    `db:"merchant_image_url" json:"merchantImageUrl"`
		OrderID           string    `db:"order_id" json:"orderId"`
		Quantity          int       `db:"quantity" json:"quantity"`
		MerchantLat       float64   `db:"merchant_lat" json:"merchantLat"`
		MerchantLon       float64   `db:"merchant_long" json:"merchantLong"`
		MerchantName      string    `db:"merchant_name" json:"merchantName"`
		MerchantCategory  string    `db:"merchant_category" json:"merchantCategory"`
		ItemName          string    `db:"item_name" json:"itemName"`
		ItemCategory      string    `db:"item_category" json:"itemCategory"`
		ItemImageURL      string    `db:"item_image_url" json:"itemImageUrl"`
		ItemPrice         int       `db:"item_price" json:"itemPrice"`
		ItemCreatedAt     time.Time `db:"item_created_at" json:"itemCreatedAt"`
		MerchantCreatedAt time.Time `db:"merchant_created_at" json:"merchantCreatedAt"`
	}

	MerchantNearbyFilter struct {
		Name             string
		UserID           string
		Offset           int
		Lat              float64
		Lon              float64
		MerchantID       string
		MerchantCategory string
		Limit            int
	}

	MerchantWithItems struct {
		Merchant
		Items []MercItem
	}
)
