package dto

import "time"

type (
	Location struct {
		Lat  float64 `json:"lat" validate:"required,min=-90,max=90"`
		Long float64 `json:"long" validate:"required,min=-180,max=180"`
	}

	Meta struct {
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	Merchant struct {
		ID        string    `json:"merchantId" db:"id"`
		Name      string    `json:"name" db:"name"`
		Category  string    `json:"merchantCategory" db:"category"`
		ImageURL  string    `json:"imageUrl" db:"image_url"`
		Location  Location  `json:"location"`
		CreatedAt time.Time `json:"createdAt" db:"created_at"`
	}

	MercItem struct {
		ID       string    `json:"itemId" db:"id"`
		Name     string    `json:"name" db:"name"`
		Category string    `json:"productCategory" db:"category"`
		ImageURL string    `json:"imageUrl" db:"image_url"`
		Price    int       `json:"price" db:"price"`
		CreateAt time.Time `json:"createdAt" db:"created_at"`
	}

	MerchantResponse struct {
		Data []Merchant `json:"data"`
		Meta Meta       `json:"meta"`
	}

	MercItemResponse struct {
		Data []MercItem `json:"data"`
		Meta Meta       `json:"meta"`
	}

	CreateMerchantRequest struct {
		Name             string   `json:"name" validate:"required,min=2,max=30"`
		MerchantCategory string   `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
		ImageURL         string   `json:"imageUrl" validate:"required,validUrl"`
		Location         Location `json:"location" validate:"required"`
	}

	CreateMercItemRequest struct {
		Name            string `json:"name" validate:"required,min=2,max=30"`
		ProductCategory string `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
		ImageURL        string `json:"imageUrl" validate:"required,validUrl"`
		Price           int    `json:"price" validate:"required,min=1"`
	}

	CreateMerchantResponse struct {
		ID string `json:"merchantId" db:"id"`
	}

	CreateMercItemResponse struct {
		ID string `json:"itemId" db:"id"`
	}
)
