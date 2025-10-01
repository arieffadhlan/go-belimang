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

	Location struct {
		Lat  float64 `json:"lat" validate:"required"`
		Long float64 `json:"long" validate:"required"`
	}

	MerchantCategory string

	MerchantsResponse struct {
		Data []Merchant `json:"data"`
		Meta Meta       `json:"meta"`
	}

	MerchantResponse struct {
		ID string `json:"merchantId" db:"id"`
	}

	MerchantCreateRequest struct {
		Name             string           `json:"name" validate:"required,min=2,max=30"`
		MerchantCategory MerchantCategory `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
		ImageURL         string           `json:"imageUrl" validate:"required"`
		Location         Location         `json:"location"`
	}
)
