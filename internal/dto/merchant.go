package dto

import "time"

type (
	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	Merchant struct {
		ID               string           `json:"merchantId" db:"id"`
		Name             string           `json:"name" db:"name"`
		MerchantCategory MerchantCategory `json:"merchantCategory" db:"merchant_category"`
		ImageURL         string           `json:"imageUrl" db:"image_url"`
		Location         Location         `json:"location" db:"location"`
		CreatedAt        time.Time        `json:"createdAt" db:"created_at"`
	}

	Item struct {
		ID              string          `json:"itemId" db:"id"`
		Name            string          `json:"name" db:"name"`
		ProductCategory ProductCategory `json:"productCategory" db:"product_category"`
		Price           int             `json:"price" db:"price"`
		ImageURL        string          `json:"imageUrl" db:"image_url"`
		CreateAt        time.Time       `json:"createdAt" db:"created_at"`
	}

	Location struct {
		Lat  float64 `json:"lat" validate:"required"`
		Long float64 `json:"long" validate:"required"`
	}

	MerchantCategory string
	ProductCategory  string

	MerchantsResponse struct {
		Data []Merchant `json:"data"`
		Meta Meta       `json:"meta"`
	}

	ItemsResponse struct {
		Data []Item `json:"data"`
		Meta Meta   `json:"meta"`
	}
	MerchantResponse struct {
		ID string `json:"merchantId" db:"id"`
	}

	ItemMerchantResponse struct {
		ID string `json:"itemId" db:"id"`
	}

	MerchantCreateRequest struct {
		Name             string           `json:"name" validate:"required,min=2,max=30"`
		MerchantCategory MerchantCategory `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
		ImageURL         string           `json:"imageUrl" validate:"required"`
		Location         Location         `json:"location"`
	}

	ItemMerchantRequest struct {
		Name            string          `json:"name" validate:"required,min=2,max=30"`
		ProductCategory ProductCategory `json:"productCategory" validate:"required, oneof=Beverage Food Snack Condiments Additions"`
		Price           int             `json:"price" validate:"required,min=1"`
		ImageURL        string          `json:"imageUrl" validate:"required"`
	}
)
