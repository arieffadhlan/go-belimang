package entities

import "time"

type (
	Merchant struct {
		ID               string           `db:"id"`
		Name             string           `db:"name"`
		MerchantCategory MerchantCategory `db:"merchant_category"`
		ImageURL         string           `db:"image_url"`
		Location         Location         `db:"location"`
		CreatedAt        time.Time        `db:"created_at"`
	}

	Location struct {
		Lat  float64
		Long float64
	}

	MerchantCategory string

	MerchantFilter struct {
		MerchantID       string
		Name             string
		MerchantCategory string
		Limit            int
		Offset           int
		SortCreatedAt    string // "asc" or "desc"
	}
)

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)
