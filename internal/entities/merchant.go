package entities

import "time"

type (
	Location struct {
		Lat  float64
		Long float64
	}

	Merchant struct {
		ID        string `db:"id"`
		Name      string `db:"name"`
		Category  string `db:"category"`
		ImageURL  string `db:"image_url"`
		Location  Location
		CreatedAt time.Time `db:"created_at"`
	}

	MerchantItem struct {
		ID         string    `db:"id"`
		Name       string    `db:"name"`
		MerchantID string    `db:"merchant_id"`
		Category   string    `db:"category"`
		ImageURL   string    `db:"image_url"`
		Price      int       `db:"price"`
		CreatedAt  time.Time `db:"created_at"`
	}

	MerchantFilter struct {
		Limit            int
		CreatedAt        string
		Name             string
		MerchantID       string
		MerchantCategory string
		Offset           int
	}

	MerchantItemFilter struct {
		Limit           int
		CreatedAt       string
		Name            string
		ItemID          string
		ProductCategory string
		Offset          int
	}
)
