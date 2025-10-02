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

	MerchantItem struct {
		ID              string          `db:"id"`
		MerchantID      string          `db:"merchant_id"`
		Name            string          `db:"name"`
		ProductCategory ProductCategory `db:"product_category"`
		Price           int             `db:"price"`
		ImageURL        string          `db:"image_url"`
		CreatedAt       time.Time       `db:"created_at"`
	}

	Location struct {
		Lat  float64
		Long float64
	}

	MerchantCategory string
	ProductCategory  string

	MerchantFilter struct {
		MerchantID       string
		Name             string
		MerchantCategory string
		Limit            int
		Offset           int
		SortCreatedAt    string // "asc" or "desc"
	}

	ItemFilter struct {
		ItemID          string
		Limit           int
		Offset          int
		Name            string
		ProductCategory string
		SortCreatedAt   string // "asc" or "desc"
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

const (
	Beverage   ProductCategory = "Beverage"
	Food       ProductCategory = "Food"
	Snack      ProductCategory = "Snack"
	Condiments ProductCategory = "Condiments"
	Additions  ProductCategory = "Additions"
)
